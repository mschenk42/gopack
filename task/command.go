package task

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

// Command runs commands via exec.Command
type Command struct {
	Name      string
	Args      []string
	Env       []string
	Dir       string
	Stream    bool
	Sensitive bool
	Timeout   time.Duration

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (c Command) Run(runActions ...action.Name) gopack.ActionRunStatus {
	c.setDefaults()
	return c.RunActions(&c, c.registerActions(), runActions)
}

func (c Command) registerActions() action.Funcs {
	return action.Funcs{
		action.Run: c.run,
	}
}

func (c *Command) setDefaults() {
	if c.Timeout == 0 {
		c.Timeout = 1 * time.Hour
	}
}

// String returns a string which identifies the task with it's property values
func (c Command) String() string {
	if c.Sensitive {
		return fmt.Sprintf("command %s %v", c.Name, Redact(c.Args...))
	}
	return fmt.Sprintf("command %s %v", c.Name, c.Args)
}

func (c Command) run() (bool, error) {
	if c.Stream {
		if err := execCmdStream(gopack.NewTaskInfoWriter(), c.Timeout, c.Name, c.Env, c.Dir, c.Args...); err != nil {
			return false, fmt.Errorf("unable to execute %s, %s", c, err)
		}
	} else {
		b, err := execCmd(c.Timeout, c.Name, c.Env, c.Dir, c.Args...)
		if err != nil {
			return false, err
		}
		if string(b) != "" {
			gopack.NewTaskInfoWriter().Write(b)
		}
	}
	return true, nil
}

func execCmd(timeout time.Duration, command string, env []string, wd string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = env
	if wd != "" {
		cmd.Dir = wd
	}
	b, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return b, ctx.Err()
	}
	if err != nil {
		return b, err
	}
	return b, nil
}

func execCmdStream(w io.Writer, timeout time.Duration, command string, env []string, wd string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Env = env
	if wd != "" {
		cmd.Dir = wd
	}
	return cmd.Run()
}

package task

import (
	"fmt"
	"time"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

// Command runs commands via exec.Command
type Command struct {
	Name    string
	Args    []string
	Stream  bool
	Timeout time.Duration

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (c Command) Run(runActions ...action.Enum) gopack.ActionRunStatus {
	c.setDefaults()
	return c.RunActions(&c, c.registerActions(), runActions)
}

func (c Command) registerActions() action.Methods {
	return action.Methods{
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
	return fmt.Sprintf("command %s %v", c.Name, c.Args)
}

func (c Command) run() (bool, error) {
	if c.Stream {
		if err := ExecCmdStream(gopack.NewTaskInfoWriter(), c.Timeout, c.Name, c.Args...); err != nil {
			return false, err
		}
	} else {
		b, err := ExecCmd(c.Timeout, c.Name, c.Args...)
		if err != nil {
			return false, err
		}
		if string(b) != "" {
			gopack.NewTaskInfoWriter().Write(b)
		}
	}
	return true, nil
}

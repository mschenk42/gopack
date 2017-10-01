package task

import (
	"fmt"
	"os/exec"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

// Package ...
type Package struct {
	Names []string

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (p Package) Run(runActions ...action.Name) gopack.ActionRunStatus {
	p.setDefaults()
	return p.RunActions(&p, p.registerActions(), runActions)
}

func (p Package) registerActions() action.Funcs {
	return action.Funcs{
		action.Install: p.install,
	}
}

func (p *Package) setDefaults() {
}

// String returns a string which identifies the task with it's property values
func (p Package) String() string {
	return fmt.Sprintf("package %+v", p.Names)
}

func (p Package) install() (bool, error) {
	c, err := p.packageCommand()
	if err != nil {
		return false, err
	}
	return c.Run(action.Run)[action.Run], nil
}

func (p Package) packageCommand() (Command, error) {
	path, err := exec.LookPath("apt-get")
	if err != nil && err.(*exec.Error).Err != exec.ErrNotFound {
		return Command{}, err
	}
	if err == nil {
		args := []string{"install", "-y"}
		args = append(args, p.Names...)
		c := Command{
			Name:   path,
			Stream: true,
			Args:   args,
		}
		return c, nil
	}

	path, err = exec.LookPath("yum")
	if err != nil && err.(*exec.Error).Err != exec.ErrNotFound {
		return Command{}, err
	}
	if err == nil {
		args := []string{"install", "-y"}
		args = append(args, p.Names...)
		c := Command{
			Name:   path,
			Stream: true,
			Args:   args,
		}
		return c, nil
	}
	return Command{}, err
}

package task

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/mschenk42/gopack"
)

type Group struct {
	Name string

	gopack.BaseTask
}

func (g Group) Run(runActions ...gopack.Action) bool {
	g.setDefaults()
	return g.RunActions(&g, g.registerActions(), runActions)
}

func (g Group) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: g.create,
		gopack.RemoveAction: g.remove,
	}
}

func (g *Group) setDefaults() {
}

func (g Group) String() string {
	return fmt.Sprintf("group %s", g.Name)
}

func (g Group) create() (bool, error) {
	if _, err := user.LookupGroup(g.Name); err == nil {
		return false, nil
	} else {
		if !strings.Contains(err.Error(), "unknown group") {
			return false, g.TaskError(g, gopack.CreateAction, err)
		}
	}
	if err := createGroupCmd(g); err != nil {
		return false, g.TaskError(g, gopack.CreateAction, err)
	}
	return true, nil
}

func (g Group) remove() (bool, error) {
	if _, err := user.LookupGroup(g.Name); err != nil {
		return false, g.TaskError(g, gopack.RemoveAction, err)
	}
	if err := removeGroupCmd(g); err != nil {
		return false, g.TaskError(g, gopack.RemoveAction, err)
	}
	return true, nil
}

func createGroupCmd(g Group) error {
	switch runtime.GOOS {
	case "linux":
		return createGroupLinux(g)
	default:
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
}

func removeGroupCmd(g Group) error {
	switch runtime.GOOS {
	case "linux":
		return removeGroupLinux(g)
	default:
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
}

func createGroupLinux(g Group) error {
	_, err := execCmd(time.Second*10, "groupadd", g.Name)
	return err
}

func removeGroupLinux(g Group) error {
	_, err := execCmd(time.Second*10, "groupdel", g.Name)
	return err
}

package task

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

type Group struct {
	Name string

	gopack.BaseTask
}

func (g Group) Run(runActions ...action.Enum) gopack.ActionRunStatus {
	g.setDefaults()
	return g.RunActions(&g, g.registerActions(), runActions)
}

func (g Group) registerActions() action.Methods {
	return action.Methods{
		action.Create: g.create,
		action.Remove: g.remove,
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
			return false, gopack.NewTaskError(g, action.Create, err)
		}
	}
	if err := createGroup(g); err != nil {
		return false, gopack.NewTaskError(g, action.Create, err)
	}
	return true, nil
}

func (g Group) remove() (bool, error) {
	if _, err := user.LookupGroup(g.Name); err != nil {
		return false, gopack.NewTaskError(g, action.Remove, err)
	}
	if err := removeGroup(g); err != nil {
		return false, gopack.NewTaskError(g, action.Remove, err)
	}
	return true, nil
}

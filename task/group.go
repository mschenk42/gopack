package task

import (
	"fmt"

	"github.com/mschenk42/gopack"
)

type Group struct {
	GroupName string

	gopack.BaseTask
}

func (u Group) Run(runActions ...gopack.Action) bool {
	u.setDefaults()
	return u.RunActions(&u, u.registerActions(), runActions)
}

func (u Group) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: u.create,
		gopack.RemoveAction: u.remove,
	}
}

func (u *Group) setDefaults() {
}

func (u Group) String() string {
	return fmt.Sprintf("group %s", u.GroupName)
}

func (u Group) create() (bool, error) {
	return true, nil
}

func (u Group) remove() (bool, error) {
	return true, nil
}

package task

import (
	"fmt"

	"github.com/mschenk42/gopack"
)

type Group struct {
	GroupName string

	props *gopack.Properties
	gopack.BaseTask
}

func (u Group) Run(props *gopack.Properties, runActions ...gopack.Action) bool {
	u.props = props
	u.setDefaults()
	return u.RunActions(&u, u.registerActions(), runActions)
}

func (u Group) Properties() *gopack.Properties {
	return u.props
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

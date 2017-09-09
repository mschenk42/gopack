package task

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

type User struct {
	Name  string
	Group string
	Home  string

	gopack.BaseTask
}

func (u User) Run(runActions ...action.Enum) gopack.ActionRunStatus {
	u.setDefaults()
	return u.RunActions(&u, u.registerActions(), runActions)
}

func (u User) registerActions() action.Methods {
	return action.Methods{
		action.Create: u.create,
		action.Remove: u.remove,
	}
}

func (u *User) setDefaults() {
}

func (u User) String() string {
	return fmt.Sprintf("user %s %s %s", u.Name, u.Group, u.Home)
}

func (u User) create() (bool, error) {
	if _, err := user.Lookup(u.Name); err == nil {
		return false, nil
	} else {
		if !strings.Contains(err.Error(), "unknown user") {
			return false, err
		}
	}
	if err := createUser(u); err != nil {
		return false, err
	}
	return true, nil
}

func (u User) remove() (bool, error) {
	if _, err := user.Lookup(u.Name); err != nil {
		return false, err
	}
	if err := removeUser(u); err != nil {
		return false, err
	}
	return true, nil
}

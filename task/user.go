package task

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/mschenk42/gopack"
)

type User struct {
	UserName string
	Group    string
	Home     string

	props *gopack.Properties
	gopack.BaseTask
}

func (u User) Run(props *gopack.Properties, runActions ...gopack.Action) bool {
	u.props = props
	u.setDefaults()
	return u.RunActions(&u, u.registerActions(), runActions)
}

func (u User) Properties() *gopack.Properties {
	return u.props
}

func (u User) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: u.create,
		gopack.RemoveAction: u.remove,
	}
}

func (u *User) setDefaults() {
}

func (u User) String() string {
	return fmt.Sprintf("user %s %s %s", u.UserName, u.Group, u.Home)
}

func (u User) create() (bool, error) {
	if _, err := user.Lookup(u.UserName); err == nil {
		return false, nil
	} else {
		if !strings.Contains(err.Error(), "unknown user") {
			return false, u.TaskError(u, gopack.CreateAction, err)
		}
	}
	switch runtime.GOOS {
	case "linux":
		return u.createUserLinux()
	default:
		return false, u.TaskError(u, gopack.CreateAction, fmt.Errorf("not supported for %s", runtime.GOOS))
	}
}

func (u User) remove() (bool, error) {
	if _, err := user.Lookup(u.UserName); err != nil {
		return false, u.TaskError(u, gopack.RemoveAction, err)
	}
	switch runtime.GOOS {
	case "linux":
		return u.removeUserLinux()
	default:
		return false, u.TaskError(u, gopack.RemoveAction, fmt.Errorf("not supported for %s", runtime.GOOS))
	}
}

func (u User) createUserLinux() (bool, error) {
	x := []string{}
	if u.Group != "" {
		x = append(x, "-g", u.Group)
	}
	if u.Home != "" {
		x = append(x, "-d", u.Home)
	}
	x = append(x, u.UserName)
	if _, err := execCmd(time.Second*10, "useradd", x...); err != nil {
		return false, u.TaskError(u, gopack.CreateAction, err)
	}
	return true, nil
}

func (u User) removeUserLinux() (bool, error) {
	if _, err := execCmd(time.Second*10, "userdel", u.UserName); err != nil {
		return false, u.TaskError(u, gopack.RemoveAction, err)
	}
	return true, nil
}

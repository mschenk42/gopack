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

	gopack.BaseTask
}

func (u User) Run(runActions ...gopack.Action) bool {
	u.setDefaults()
	return u.RunActions(&u, u.registerActions(), runActions)
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
	if err := createUserCmd(u); err != nil {
		return false, u.TaskError(u, gopack.CreateAction, err)
	}
	return true, nil
}

func (u User) remove() (bool, error) {
	if _, err := user.Lookup(u.UserName); err != nil {
		return false, u.TaskError(u, gopack.RemoveAction, err)
	}
	if err := removeUserCmd(u); err != nil {
		return false, u.TaskError(u, gopack.RemoveAction, err)
	}
	return true, nil
}

func createUserCmd(u User) error {
	switch runtime.GOOS {
	case "linux":
		return createUserLinux(u)
	default:
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
}

func removeUserCmd(u User) error {
	switch runtime.GOOS {
	case "linux":
		return removeUserLinux(u)
	default:
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
}

func createUserLinux(u User) error {
	x := []string{}
	if u.Group != "" {
		x = append(x, "-g", u.Group)
	}
	if u.Home != "" {
		x = append(x, "-d", u.Home)
	}
	x = append(x, u.UserName)
	if _, err := execCmd(time.Second*10, "useradd", x...); err != nil {
		return err
	}
	return nil
}

func removeUserLinux(u User) error {
	_, err := execCmd(time.Second*10, "userdel", u.UserName)
	return err
}

package task

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
	"time"

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
			return false, gopack.NewTaskError(u, action.Create, err)
		}
	}
	if err := createUserCmd(u); err != nil {
		return false, gopack.NewTaskError(u, action.Create, err)
	}
	return true, nil
}

func (u User) remove() (bool, error) {
	if _, err := user.Lookup(u.Name); err != nil {
		return false, gopack.NewTaskError(u, action.Remove, err)
	}
	if err := removeUserCmd(u); err != nil {
		return false, gopack.NewTaskError(u, action.Remove, err)
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
	x = append(x, u.Name)
	if _, err := execCmd(time.Second*10, "useradd", x...); err != nil {
		return err
	}
	return nil
}

func removeUserLinux(u User) error {
	_, err := execCmd(time.Second*10, "userdel", u.Name)
	return err
}

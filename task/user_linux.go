package task

import (
	"strings"
	"time"

	"github.com/mschenk42/gopack/action"
)

func createUser(u User) {
	args := []string{}
	if u.Group != "" {
		args = append(args, "-g", u.Group)
	}
	if len(u.Groups) > 0 {
		args = append(args, "-G", strings.Join(u.Groups, ","))
	}
	if u.Home != "" {
		args = append(args, "-d", u.Home)
	}
	args = append(args, u.Name)
	Command{
		Name:    "useradd",
		Args:    args,
		Timeout: time.Second * 10,
	}.Run(action.Run)
}

func removeUser(u User) {
	Command{
		Name:    "userdel",
		Args:    []string{u.Name},
		Timeout: time.Second * 10,
	}.Run(action.Run)
}

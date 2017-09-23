package task

import (
	"time"

	"github.com/mschenk42/gopack/action"
)

func createGroup(g Group) {
	Command{
		Name:    "groupadd",
		Args:    []string{g.Name},
		Timeout: time.Second * 10,
	}.Run(action.Run)
}

func removeGroup(g Group) {
	Command{
		Name:    "groupdel",
		Args:    []string{g.Name},
		Timeout: time.Second * 10,
	}.Run(action.Run)
}

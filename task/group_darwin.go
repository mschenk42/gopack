package task

import (
	"fmt"
	"runtime"
)

func createGroup(g Group) error {
	panic(fmt.Sprintf("create group not implemented for %s", runtime.GOOS))
}

func removeGroup(g Group) error {
	panic(fmt.Sprintf("remove group not implemented for %s", runtime.GOOS))
}

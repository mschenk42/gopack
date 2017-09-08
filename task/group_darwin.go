package task

import (
	"fmt"
	"runtime"
)

func createGroup(g Group) error {
	return fmt.Errorf("create group not implemented for %s", runtime.GOOS)
}

func removeGroup(g Group) error {
	return fmt.Errorf("remove group not implemented for %s", runtime.GOOS)
}

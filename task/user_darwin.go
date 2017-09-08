package task

import (
	"fmt"
	"runtime"
)

func createUser(u User) error {
	return fmt.Errorf("create user not implemented for %s", runtime.GOOS)
}

func removeUser(u User) error {
	return fmt.Errorf("remove user not implemented for %s", runtime.GOOS)
}

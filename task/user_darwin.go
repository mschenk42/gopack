package task

import (
	"fmt"
	"runtime"
)

func createUser(u User) {
	panic(fmt.Sprintf("create user not implemented for %s", runtime.GOOS))
}

func removeUser(u User) {
	panic(fmt.Sprintf("remove user not implemented for %s", runtime.GOOS))
}

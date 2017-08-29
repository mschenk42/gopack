// +build linux

package task

import (
	"testing"

	"github.com/mschenk42/gopack"
)

func TestCreateUser(t *testing.T) {
	User{
		UserName: "test",
	}.Run(gopack.CreateAction)
}

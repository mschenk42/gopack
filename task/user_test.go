package task

import (
	"bytes"
	"fmt"
	"log"
	"os/user"
	"runtime"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux only test")
	}
	assert := assert.New(t)

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	x := User{
		Name: "test",
	}

	defer func() {
		assert.NotPanics(func() { x.remove() }, "x.remove() %s", x)
	}()

	assert.NotPanics(func() { x.Run(action.Create) }, "action.Create %s", x)
	assert.Regexp(`.*user test.*create.*(run)`, buf.String())
	_, err := user.Lookup(x.Name)
	assert.NoError(err)
	fmt.Print(buf.String())
}

func TestRemoveUserLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux only test")
	}
	assert := assert.New(t)

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	x := User{
		Name: "test",
	}

	defer func() {
		assert.NotPanics(func() { x.remove() }, "x.remove() %s", x)
	}()

	assert.NotPanics(func() { x.Run(action.Create) }, "action.Create %s", x)
	assert.NotPanics(func() { x.Run(action.Remove) }, "action.Remove %s", x)
	assert.Regexp(`.*user test.*remove.*(run)`, buf.String())
	_, err := user.Lookup(x.Name)
	assert.NotNil(err)
	fmt.Print(buf.String())
}

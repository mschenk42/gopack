package task

import (
	"bytes"
	"fmt"
	"log"
	"os/user"
	"runtime"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroupLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux only test")
	}
	assert := assert.New(t)

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	x := Group{
		Name: "test",
	}

	defer func() {
		x.remove()
	}()

	assert.NotPanics(func() { x.Run(gopack.CreateAction) })
	assert.Regexp(`.*group test.*create.*(run)`, buf.String())
	_, err := user.LookupGroup(x.Name)
	assert.NoError(err)
	fmt.Print(buf.String())
}

func TestRemoveGroupLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux only test")
	}
	assert := assert.New(t)

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	x := Group{
		Name: "test",
	}

	defer func() {
		x.remove()
	}()

	assert.NotPanics(func() { x.Run(gopack.CreateAction) })
	assert.NotPanics(func() { x.Run(gopack.RemoveAction) })
	assert.Regexp(`.*group test.*remove.*(run)`, buf.String())
	_, err := user.LookupGroup(x.Name)
	assert.NotNil(err)
	fmt.Print(buf.String())
}

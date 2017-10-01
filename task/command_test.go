package task

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/stretchr/testify/assert"
)

func TestCommandStream(t *testing.T) {
	assert := assert.New(t)

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	c := Command{
		Name:   "echo",
		Args:   []string{"hello"},
		Stream: true,
	}

	assert.Equal(gopack.ActionRunStatus{action.Run: true}, c.Run(action.Run))
	assert.Regexp(".*hello\n", buf.String())
	assert.Regexp(`.*command.*echo.*hello.*run.*(started)`, buf.String())
	assert.Regexp(`.*command.*echo.*hello.*run.*(has run)`, buf.String())
	fmt.Print(buf.String())
}

func TestExecCmdStreamFunc(t *testing.T) {
	assert := assert.New(t)

	buf := &bytes.Buffer{}
	assert.NoError(execCmdStream(buf, 1*time.Second, "echo", []string{}, "", "hello"))
	assert.Equal("hello\n", buf.String())
}

func TestExecCmdFunc(t *testing.T) {
	assert := assert.New(t)

	b, err := execCmd(1*time.Second, "echo", []string{}, "", "hello")
	assert.NoError(err)
	assert.Equal("hello\n", string(b))
}

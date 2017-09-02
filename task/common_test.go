package task

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecCmdStreamFunc(t *testing.T) {
	assert := assert.New(t)

	buf := &bytes.Buffer{}
	assert.NoError(execCmdStream(buf, 1*time.Second, "echo", "hello"))
	assert.Equal("hello\n", buf.String())
}

func TestExecCmdFunc(t *testing.T) {
	assert := assert.New(t)

	b, err := execCmd(1*time.Second, "echo", "hello")
	assert.NoError(err)
	assert.Equal("hello\n", string(b))
}

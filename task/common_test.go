package task

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandStream(t *testing.T) {
	assert := assert.New(t)

	buf := &bytes.Buffer{}
	assert.NoError(execCmdStream(buf, 10*time.Second, "echo", "hello"))
	assert.Equal("hello", buf.String())
}

// +build linux

package mypack

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux only test")
	}
	assert := assert.New(t)
	assert.NotPanics(func() { Run(nil) })
}

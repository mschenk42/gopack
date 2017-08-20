package mypack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() { Run(nil) })
}

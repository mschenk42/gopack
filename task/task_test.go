package task

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newLoggerBuffer() (*log.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	return logger, buf
}

func TestRunTask(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
		Base: Base{NotIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Not-Run", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		Base: Base{NotIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		Base: Base{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		Base: Base{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Not-Run", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		Base: Base{
			OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil },
			NotIf:  func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
		Base: Base{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil, logger) })
	assert.Regexp("unable to run task1", buf.String())
	fmt.Print(buf.String())
}

type Task1 struct {
	Name string
	Base
}

func (t Task1) Run(props Properties, logger *log.Logger, runActions ...Action) bool {
	regActions := ActionMethods{
		Create: t.create,
	}
	return t.Base.RunActions(&t, regActions, runActions, props, logger)
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create(props Properties, logger *log.Logger) (bool, error) {
	return true, nil
}

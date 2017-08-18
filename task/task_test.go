package task

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)

	logInfoSave := LogInfo
	logErrSave := LogErr
	defer func() {
		LogInfo = logInfoSave
		LogErr = logErrSave
	}()

	buf := &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)

	logInfoSave := LogInfo
	logErrSave := LogErr
	defer func() {
		LogInfo = logInfoSave
		LogErr = logErrSave
	}()

	buf := &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)
	t1 := Task1{
		Name: "task1",
		Base: Base{NotIf: func(p Properties) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Not-Run", buf.String())
	fmt.Print(buf.String())

	buf = &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)
	t1 = Task1{
		Name: "task1",
		Base: Base{NotIf: func(p Properties) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())

	buf = &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)
	t1 = Task1{
		Name: "task1",
		Base: Base{OnlyIf: func(p Properties) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())

	buf = &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)
	t1 = Task1{
		Name: "task1",
		Base: Base{OnlyIf: func(p Properties) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Not-Run", buf.String())
	fmt.Print(buf.String())

	buf = &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)
	t1 = Task1{
		Name: "task1",
		Base: Base{
			OnlyIf: func(p Properties) (bool, error) { return true, nil },
			NotIf:  func(p Properties) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, Create) })
	assert.Regexp("task1.*create.*Did-Run", buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)

	logInfoSave := LogInfo
	logErrSave := LogErr
	defer func() {
		LogInfo = logInfoSave
		LogErr = logErrSave
	}()

	buf := &bytes.Buffer{}
	LogInfo = log.New(buf, "INFO: ", 0)
	LogErr = log.New(buf, "ERROR: ", 0)

	t1 := Task1{
		Name: "task1",
		Base: Base{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil) })
	assert.Regexp("ERROR.*Unable to run task1", buf.String())
	fmt.Print(buf.String())
}

type Task1 struct {
	Name string
	Base
}

func (t Task1) Run(props Properties, runActions ...Action) bool {
	regActions := ActionMethods{
		Create: t.create,
	}
	return t.Base.RunActions(&t, regActions, runActions, props)
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create(props Properties) (bool, error) {
	return true, nil
}

type Task2 struct {
	Name string
	Base
}

func (t Task2) Run(props Properties, runActions ...Action) bool {
	regActions := ActionMethods{
		Nothing: t.nothing,
		Create:  t.create,
	}
	return t.Base.RunActions(&t, regActions, runActions, props)
}

func (t Task2) String() string {
	return t.Name
}

func (t Task2) create(props Properties) (bool, error) {
	return true, nil
}

func (t Task2) nothing(props Properties) (bool, error) {
	return false, nil
}

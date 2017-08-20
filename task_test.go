package gopack

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[RUN\\]", buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[NOT RUN\\]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[RUN\\]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[RUN\\]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[NOT RUN\\]", buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		BaseTask: BaseTask{
			OnlyIf: func(p Properties, x *log.Logger) (bool, error) { return true, nil },
			NotIf:  func(p Properties, x *log.Logger) (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[RUN\\]", buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil, logger) })
	assert.Regexp("unable to run task1", buf.String())
	fmt.Print(buf.String())
}

func TestNotify(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	t2 := Task2{
		Name: "task2 notified",
	}

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			Notify: map[Action][]func(){
				CreateAction: []func(){
					func() { t2.Run(nil, logger, CreateAction) },
				},
			},
		},
	}

	assert.NotPanics(func() { t1.Run(nil, logger, CreateAction) })
	assert.Regexp("task1.*create.*\\[RUN\\]", buf.String())
	assert.Regexp("task2 notified.*create.*\\[RUN\\]", buf.String())
	fmt.Print(buf.String())
}

func newLoggerBuffer() (*log.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	return logger, buf
}

type Task1 struct {
	Name string
	BaseTask
}

func (t Task1) Run(props Properties, logger *log.Logger, runActions ...Action) bool {
	regActions := ActionMethods{
		CreateAction: t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions, props, logger)
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create(props Properties, logger *log.Logger) (bool, error) {
	return true, nil
}

type Task2 struct {
	Name string
	BaseTask
}

func (t Task2) Run(props Properties, logger *log.Logger, runActions ...Action) bool {
	regActions := ActionMethods{
		NothingAction: t.nothing,
		CreateAction:  t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions, props, logger)
}

func (t Task2) String() string {
	return t.Name
}

func (t Task2) create(props Properties, logger *log.Logger) (bool, error) {
	return true, nil
}

func (t Task2) nothing(props Properties, logger *log.Logger) (bool, error) {
	return false, nil
}

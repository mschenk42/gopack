package gopack

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	failKeywords = "\\(up to date\\)"
	passKeywords = "\\(run\\) "
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", "skipped due to not_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", "run due to not_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", "run due to only_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", "skipped due to only_if"), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		BaseTask: BaseTask{
			OnlyIf: func() (bool, error) { return true, nil },
			NotIf:  func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", "skipped due to not_if"), buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil) })
	assert.Regexp("unable to run task1", buf.String())
	fmt.Print(buf.String())
}

func TestWhenRun(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t1.NotifyWhen(t2, CreateAction, nil, false)

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf("task2 notified.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestDelayedRun(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Log
	buf := &bytes.Buffer{}
	Log = log.New(buf, "", 0)
	defer func() { Log = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t3 := Task1{
		Name: "task3",
	}

	t1.NotifyWhen(t2, CreateAction, nil, true)

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.NotPanics(func() { t3.Run(nil, CreateAction) })
	assert.NotPanics(func() { DelayedNotify.Run() })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf("task3.*create.*%s", passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf("task2 notified.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())
}

type Task1 struct {
	Name string

	logger     *log.Logger
	properties *Properties
	BaseTask
}

func (t Task1) Run(props *Properties, runActions ...Action) bool {
	t.properties = props
	regActions := ActionMethods{
		CreateAction: t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions)
}

func (d Task1) Properties() *Properties {
	return d.properties
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create() (bool, error) {
	return true, nil
}

type Task2 struct {
	Name string

	logger     *log.Logger
	properties *Properties
	BaseTask
}

func (t Task2) Run(props *Properties, runActions ...Action) bool {
	t.properties = props
	regActions := ActionMethods{
		NothingAction: t.nothing,
		CreateAction:  t.create,
	}
	return t.BaseTask.RunActions(&t, regActions, runActions)
}

func (d Task2) Properties() *Properties {
	return d.properties
}

func (t Task2) String() string {
	return t.Name
}

func (t Task2) create() (bool, error) {
	return true, nil
}

func (t Task2) nothing() (bool, error) {
	return false, nil
}

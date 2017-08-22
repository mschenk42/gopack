package gopack

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	failKeywords = "\\[NOT RUN\\]"
	passKeywords = "\\[RUN\\] "
)

func TestRunTask(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Logger
	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)
	defer func() { Logger = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestGuards(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Logger
	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)
	defer func() { Logger = saveLogger }()

	t1 := Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", failKeywords), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{NotIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name:     "task1",
		BaseTask: BaseTask{OnlyIf: func() (bool, error) { return false, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", failKeywords), buf.String())
	fmt.Print(buf.String())

	buf.Reset()
	t1 = Task1{
		Name: "task1",
		BaseTask: BaseTask{
			OnlyIf: func() (bool, error) { return true, nil },
			NotIf:  func() (bool, error) { return true, nil }},
	}
	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", failKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestContOnError(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Logger
	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)
	defer func() { Logger = saveLogger }()

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			ContOnError: true},
	}

	assert.NotPanics(func() { t1.Run(nil) })
	assert.Regexp("unable to run task1", buf.String())
	fmt.Print(buf.String())
}

func TestSubscriberStruct(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Logger
	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)
	defer func() { Logger = saveLogger }()

	t2 := Task2{
		Name: "task2 notified",
	}

	t1 := Task1{
		Name: "task1",
		BaseTask: BaseTask{
			Subscribers: map[Action][]func(){
				CreateAction: []func(){
					func() { t2.Run(nil, CreateAction) },
				},
			},
		},
	}

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
	assert.Regexp(fmt.Sprintf("task2 notified.*create.*%s", passKeywords), buf.String())
	fmt.Print(buf.String())
}

func TestSubscriberMethod(t *testing.T) {
	assert := assert.New(t)

	saveLogger := Logger
	buf := &bytes.Buffer{}
	Logger = log.New(buf, "", 0)
	defer func() { Logger = saveLogger }()

	t1 := Task1{
		Name: "task1",
	}

	t2 := Task2{
		Name: "task2 notified",
	}

	t1.AddSubscriber(t2, CreateAction, nil)

	assert.NotPanics(func() { t1.Run(nil, CreateAction) })
	assert.Regexp(fmt.Sprintf("task1.*create.*%s", passKeywords), buf.String())
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

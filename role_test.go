package gopack

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

func TestCreateRole(t *testing.T) {
	assert := assert.New(t)

	role := Role{
		Name: "role1",
	}

	assert.Equal(role.Name, "role1")
}

func TestRegisterTask(t *testing.T) {
	assert := assert.New(t)

	r := &Role{
		Name: "role1",
	}

	r.Register(
		Task1{
			Name: "task1",
		},
		CreateAction,
	)

	assert.Equal(len(r.tasks), 1)
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	r := &Role{
		Name: "role1",
	}

	r.Register(
		Task1{
			Name: "task1",
		},
		CreateAction,
	)

	assert.NotPanics(func() { r.Run(nil) })
	assert.Equal(len(r.tasks), 1)
}

func TestMergeProps(t *testing.T) {
	assert := assert.New(t)

	r := &Role{
		Name: "role1",
		Props: Properties{
			"role.prop1": "prop1",
			"role.prop2": "prop2",
		},
	}

	r.Register(
		Task1{
			Name: "task1",
		},
		CreateAction,
	)

	p := Properties{
		"role.prop2": "updated2",
		"role.prop3": "prop3",
	}

	assert.NotPanics(func() { r.Run(p) })
	assert.EqualValues(
		r.Props,
		Properties{
			"role.prop3": "prop3",
			"role.prop1": "prop1",
			"role.prop2": "updated2",
		},
	)
}

func TestDelayedRun(t *testing.T) {
	assert := assert.New(t)
	logger, buf := newLoggerBuffer()

	r := &Role{
		Name:   "role1",
		Logger: logger,
	}

	t1 := Task1{
		Name: "task1",
	}

	r.Register(t1, CreateAction)

	t2 := Task2{
		Name: "task2",
	}

	r.Register(t2, NothingAction)
	r.DelayRun(t2, t1, CreateAction)

	assert.NotPanics(func() { r.Run(nil) })
	assert.Equal(len(r.tasks), 2)
	assert.Regexp("task1.*create.*[NOT RUN]", buf.String())
	assert.Regexp("task2.*nothing.*[RUN]", buf.String())
	assert.Regexp("task2.*create.*[NOT RUN]", buf.String())
	fmt.Print(buf.String())
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

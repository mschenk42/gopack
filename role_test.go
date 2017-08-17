package mincfg

import (
	"testing"

	"github.com/mschenk42/mincfg/task"
	"github.com/stretchr/testify/assert"
)

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
		task.Create,
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
		task.Create,
	)

	assert.NotPanics(func() { r.Run(nil) })
	assert.Equal(len(r.tasks), 1)
}

func TestMergeProps(t *testing.T) {
	assert := assert.New(t)

	r := &Role{
		Name: "role1",
		Props: task.Properties{
			"role.prop1": "prop1",
			"role.prop2": "prop2",
		},
	}

	r.Register(
		Task1{
			Name: "task1",
		},
		task.Create,
	)

	p := task.Properties{
		"role.prop2": "updated2",
		"role.prop3": "prop3",
	}

	assert.NotPanics(func() { r.Run(p) })
	assert.EqualValues(
		r.Props,
		task.Properties{
			"role.prop3": "prop3",
			"role.prop1": "prop1",
			"role.prop2": "updated2",
		},
	)
}

func TestDelayedRun(t *testing.T) {
	assert := assert.New(t)

	r := &Role{
		Name: "role1",
	}

	t1 := Task1{
		Name: "task1",
	}

	r.Register(t1, task.Create)

	t2 := Task2{
		Name: "task2",
	}

	r.Register(t2, task.Nothing)
	r.DelayRun(t2, t1, task.Create)

	assert.NotPanics(func() { r.Run(nil) })
	assert.Equal(len(r.tasks), 2)
}

type Task1 struct {
	Name string
	task.Base
}

func (t Task1) Run(props task.Properties, runActions ...task.Action) bool {
	regActions := task.ActionMethods{
		task.Create: t.create,
	}
	return t.Base.RunActions(&t, regActions, runActions, props)
}

func (t Task1) String() string {
	return t.Name
}

func (t Task1) create(props task.Properties) (bool, error) {
	return true, nil
}

type Task2 struct {
	Name string
	task.Base
}

func (t Task2) Run(props task.Properties, runActions ...task.Action) bool {
	regActions := task.ActionMethods{
		task.Nothing: t.nothing,
		task.Create:  t.create,
	}
	return t.Base.RunActions(&t, regActions, runActions, props)
}

func (t Task2) String() string {
	return t.Name
}

func (t Task2) create(props task.Properties) (bool, error) {
	return true, nil
}

func (t Task2) nothing(props task.Properties) (bool, error) {
	return false, nil
}

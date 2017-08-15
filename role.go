package mincfg

import "github.com/mschenk42/mincfg/task"

type Role struct {
	Name  string
	Props task.Properties
	tasks []ActionSet
}

type ActionSet struct {
	Task    task.Task
	Actions []task.Action
}

func (r *Role) Register(t task.Task, runActions ...task.Action) {
	if r.tasks == nil {
		r.tasks = []ActionSet{}
	}
	r.tasks = append(r.tasks, ActionSet{Task: t, Actions: runActions})
}

func (r *Role) Run(props task.Properties) {
	r.Props.Merge(props)
	for _, t := range r.tasks {
		t.Task.Run(r.Props, t.Actions...)
	}
}

func (r *Role) Notify() {
	for _, t := range r.tasks {
		t.Task.Notify()
	}
}

func (r *Role) handleError(err error) {
	if err != nil {
		panic(err)
	}
}

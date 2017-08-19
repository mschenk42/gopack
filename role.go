package gopack

import (
	"log"
	"os"
)

type Role struct {
	Name         string
	Props        Properties
	Logger       *log.Logger
	tasks        []taskActions
	tasksDelayed []taskRunWhen
}

type taskRunWhen struct {
	runTask  Task
	whenTask Task
	action   Action
}

type taskActions struct {
	task    Task
	actions []Action
}

func (r *Role) Register(t Task, runActions ...Action) {
	if r.tasks == nil {
		r.tasks = []taskActions{}
	}
	r.tasks = append(r.tasks, taskActions{task: t, actions: runActions})
}

func (r *Role) DelayRun(runTask, whenTask Task, action Action) {
	r.tasksDelayed = append(
		r.tasksDelayed,
		taskRunWhen{
			runTask:  runTask,
			whenTask: whenTask,
			action:   action,
		},
	)
}

func (r *Role) Run(props Properties) {
	if r.Logger == nil {
		r.Logger = log.New(os.Stdout, "", 0)
	}

	tasksRun := []taskActions{}
	r.Props.Merge(props)
	for _, ta := range r.tasks {
		if ta.task.Run(r.Props, r.Logger, ta.actions...) {
			tasksRun = append(tasksRun, ta)
		}
	}

	tasksDelayRunned := map[string]bool{}
	for _, x := range r.tasksDelayed {
		if hasRun, _ := tasksDelayRunned[x.runTask.String()]; hasRun {
			continue
		}
		for _, y := range tasksRun {
			if y.task.String() == x.whenTask.String() {
				x.runTask.Run(props, r.Logger, x.action)
				tasksDelayRunned[x.runTask.String()] = true
			}
		}
	}
}

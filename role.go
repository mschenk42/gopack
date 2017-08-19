package mincfg

import (
	"log"
	"os"

	"github.com/mschenk42/mincfg/task"
)

type Role struct {
	Name         string
	Props        task.Properties
	tasks        []taskActions
	tasksDelayed []taskRunWhen
}

type taskRunWhen struct {
	runTask  task.Task
	whenTask task.Task
	action   task.Action
}

type taskActions struct {
	task    task.Task
	actions []task.Action
}

func (r *Role) Register(t task.Task, runActions ...task.Action) {
	if r.tasks == nil {
		r.tasks = []taskActions{}
	}
	r.tasks = append(r.tasks, taskActions{task: t, actions: runActions})
}

func (r *Role) DelayRun(runTask, whenTask task.Task, action task.Action) {
	r.tasksDelayed = append(
		r.tasksDelayed,
		taskRunWhen{
			runTask:  runTask,
			whenTask: whenTask,
			action:   action,
		},
	)
}

func (r *Role) Run(props task.Properties, logger *log.Logger) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	tasksRun := []taskActions{}
	r.Props.Merge(props)
	for _, ta := range r.tasks {
		if ta.task.Run(r.Props, logger, ta.actions...) {
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
				x.runTask.Run(props, logger, x.action)
				tasksDelayRunned[x.runTask.String()] = true
			}
		}
	}
}

package gopack

import (
	"fmt"
	"time"

	"github.com/mschenk42/gopack/action"
)

var DelayedNotify taskRunSet = taskRunSet{}

type GuardFunc func() (bool, error)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool

	notify actionTaskRunSet
}

type actionTaskRunSet map[action.Enum]map[string]func()

type taskRunSet map[string]func()

func (d *taskRunSet) Run() {
	for _, f := range *d {
		f()
	}
	//clear the list
	d = &taskRunSet{}
}

type Runner interface {
	Run(actions ...action.Enum) bool
}

type Task interface {
	Runner
	fmt.Stringer
}

func (b BaseTask) RunActions(task Task, regActions action.Methods, runActions []action.Enum) bool {
	if len(runActions) == 0 {
		b.logRunStatus(false, false, "error", task, action.Nil, time.Now())
		b.handleError(fmt.Errorf("unable to run, no action given"))
		return false
	}

	t := time.Now()
	hasRun := false
	canRun, reason := b.canRun()
	if !canRun {
		b.logRunStatus(hasRun, canRun, reason, task, runActions[0], t)
		return hasRun
	}

	for _, a := range runActions {
		f, found := regActions.Method(a)
		if !found {
			b.handleError(b.TaskError(task, a, action.ErrActionNotRegistered))
			return hasRun
		}
		hasRun, err := f()
		b.handleError(err)
		b.logRunStatus(hasRun, canRun, reason, task, a, t)
		b.notifyTasks(a)
	}

	return hasRun
}

func (b *BaseTask) NotifyWhen(notify Task, forAction, whenAction action.Enum, delayed bool) {
	if b.notify == nil {
		b.notify = actionTaskRunSet{}
	}
	if b.notify[whenAction] == nil {
		b.notify[whenAction] = map[string]func(){}
	}
	if delayed {
		DelayedNotify[fmt.Sprintf("%s:%s", notify, forAction)] = func() { notify.Run(forAction) }
	} else {
		b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
			notify.Run(forAction)
		}
	}
}

func (b BaseTask) notifyTasks(action action.Enum) {
	funcs, found := b.notify[action]
	if found {
		for _, f := range funcs {
			f()
		}
	}
}

func (b BaseTask) canRun() (bool, string) {
	var (
		err    error
		run    bool = true
		reason string
	)
	if b.OnlyIf != nil {
		reason = "due to only_if"
		run, err = b.OnlyIf()
		b.handleError(err)
	}
	if b.NotIf != nil {
		reason = "due to not_if"
		run, err = b.NotIf()
		run = !run
		b.handleError(err)
	}
	return run, reason
}

const (
	TaskLogInfoFormat  = "  * %s: %s (%s) %s"
	TaskLogErrorFormat = "    ! %s"
)

func (b BaseTask) logRunStatus(hasRun, canRun bool, reason string, task Task, action action.Enum, startTime time.Time) {
	status := ""
	switch {
	case !canRun && reason != "":
		status = fmt.Sprintf("skipped %s", reason)
	case !canRun:
		status = fmt.Sprintf("skipped", reason)
	case hasRun && reason != "":
		status = fmt.Sprintf("run %s", reason)
	case hasRun:
		status = "run"
	default:
		status = "up to date"
	}

	Log.Printf(TaskLogInfoFormat, task, action, status, time.Since(startTime))
}

func (b BaseTask) TaskError(task fmt.Stringer, action action.Enum, err error) error {
	if err == nil {
		return nil
	}
	Log.Printf(TaskLogInfoFormat, task, action, "error", time.Since(time.Now()))
	return err
}

func (b BaseTask) handleError(err error) {
	if err == nil {
		return
	}
	if b.ContOnError {
		Log.Printf(TaskLogErrorFormat, err)
	} else {
		Log.Fatalf(TaskLogErrorFormat, err)
	}
}

package gopack

import (
	"fmt"
	"io"
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
		LogTaskStatus(false, false, false, "error", task, action.Nil, time.Now())
		handleTaskError(fmt.Errorf("unable to run, no action given"), b.ContOnError)
		return false
	}

	var (
		hasRun bool
		err    error
	)

	t := time.Now()
	canRun, reason := b.canRun()
	if !canRun {
		LogTaskStatus(false, hasRun, canRun, reason, task, runActions[0], t)
		return hasRun
	}

	for _, a := range runActions {
		f, found := regActions.Method(a)
		if !found {
			handleTaskError(NewTaskError(task, a, action.ErrActionNotRegistered), b.ContOnError)
			return hasRun
		}
		hasRun, err = f()
		handleTaskError(err, b.ContOnError)
		LogTaskStatus(false, hasRun, canRun, reason, task, a, t)
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
		handleTaskError(err, b.ContOnError)
	}
	if b.NotIf != nil {
		reason = "due to not_if"
		run, err = b.NotIf()
		run = !run
		handleTaskError(err, b.ContOnError)
	}
	return run, reason
}

const (
	LogHeaderFormat = "  * %s: %s (%s) %s"
	LogErrorFormat  = "    ! %s"
	LogInfoFormat   = "    %s"
)

func NewTaskInfoWriter() io.Writer {
	return taskInfoWriter{}
}

type taskInfoWriter struct {
}

func (t taskInfoWriter) Write(b []byte) (int, error) {
	Log.Printf(LogInfoFormat, b)
	return len(b), nil
}

func LogTaskStatus(isRunning, hasRun, canRun bool, reason string, task Task, action action.Enum, startTime time.Time) {
	status := ""
	switch {
	case isRunning:
		status = "running"
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
	Log.Printf(LogHeaderFormat, task, action, status, time.Since(startTime))
}

func NewTaskError(task fmt.Stringer, action action.Enum, err error) error {
	if err == nil {
		return nil
	}
	Log.Printf(LogHeaderFormat, task, action, "error", time.Since(time.Now()))
	return err
}

func handleTaskError(err error, contOnError bool) {
	if err == nil {
		return
	}
	if contOnError {
		Log.Printf(LogErrorFormat, err)
	} else {
		Log.Fatalf(LogErrorFormat, err)
	}
}

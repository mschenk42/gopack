package gopack

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mschenk42/gopack/action"
)

var (
	DelayedNotify          taskRunSet = taskRunSet{}
	ErrActionNotRegistered            = errors.New("action not registered with task")
)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool

	notify actionTaskRunSet
}

type GuardFunc func() (bool, error)
type ActionRunStatus map[action.Enum]bool
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
	Run(actions ...action.Enum) ActionRunStatus
}

type Task interface {
	Runner
	fmt.Stringer
}

func (b BaseTask) RunActions(task Task, regActions action.Methods, runActions []action.Enum) ActionRunStatus {
	t := time.Now()

	if len(runActions) == 0 {
		TaskStatus{Task: task, Actions: action.NewSlice(action.Nil), Reason: "error", StartedAt: t}.Log()
		handleTaskError(fmt.Errorf("unable to run, no action given"), b.ContOnError)
		return ActionRunStatus{}
	}

	runStatus := ActionRunStatus{}
	for _, a := range runActions {
		runStatus[a] = false
	}

	// can we run the actions for this command?
	canRun, reason := b.canRun()
	if !canRun {
		TaskStatus{Task: task, Actions: runActions, Reason: reason, StartedAt: t}.Log()
		return runStatus
	}

	for _, a := range runActions {
		f, found := regActions.Method(a)
		if !found {
			TaskStatus{Task: task, Actions: action.NewSlice(a), HasRun: false, CanRun: true, Reason: "error", StartedAt: t}.Log()
			handleTaskError(NewTaskError(task, a, ErrActionNotRegistered), b.ContOnError)
			runStatus[a] = false
			continue
		}
		hasRun, err := f()
		if err == nil {
			TaskStatus{Task: task, Actions: action.NewSlice(a), HasRun: hasRun, CanRun: true, Reason: reason, StartedAt: t}.Log()
		}
		handleTaskError(err, b.ContOnError)
		b.notifyTasks(a)
		runStatus[a] = hasRun
	}
	return runStatus
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
	logHeaderFormat = "* %s: %s (%s) %s"
	logErrorFormat  = "! %s"
	logInfoFormat   = "%s"
)

var colorize = ColorFormat{}

func NewTaskInfoWriter() io.Writer {
	return taskInfoWriter{}
}

type taskInfoWriter struct {
}

func (t taskInfoWriter) Write(b []byte) (int, error) {
	Log.Printf(logInfoFormat, b)
	return len(b), nil
}

type TaskStatus struct {
	Task      Task
	Actions   []action.Enum
	IsRunning bool
	HasRun    bool
	CanRun    bool
	Reason    string
	StartedAt time.Time
}

func (t TaskStatus) Log() {
	status := ""
	switch {
	case t.IsRunning:
		status = "running"
	case !t.CanRun && t.Reason != "":
		status = fmt.Sprintf("skipped %s", t.Reason)
	case !t.CanRun:
		status = "skipped"
	case t.HasRun && t.Reason != "":
		status = fmt.Sprintf("run %s", t.Reason)
	case t.HasRun:
		status = "run"
	default:
		status = "up to date"
	}
	actions := []string{}
	for _, a := range t.Actions {
		actions = append(actions, a.String())

	}
	colorFunc := colorize.Green
	if t.Reason == "error" {
		colorFunc = colorize.Red
	}
	Log.Printf(colorFunc(logHeaderFormat), t.Task, strings.Join(actions, ","), status, time.Since(t.StartedAt))
}

func NewTaskError(task fmt.Stringer, action action.Enum, err error) error {
	if err == nil {
		return nil
	}
	Log.Printf(colorize.Red(logHeaderFormat), task, action, "error", time.Since(time.Now()))
	return err
}

func handleTaskError(err error, contOnError bool) {
	if err == nil {
		return
	}
	if contOnError {
		Log.Printf(colorize.Red(logErrorFormat), err)
	} else {
		Log.Fatalf(colorize.Red(logErrorFormat), err)
	}
}

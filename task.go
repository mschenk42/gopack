package gopack

import (
	"errors"
	"fmt"
	"time"
)

const (
	AddAction Action = iota
	CreateAction
	DisableAction
	EnableAction
	InstallAction
	LockAction
	NilAction
	NothingAction
	RemoveAction
	RunAction
	ReloadAction
	RestartAction
	StartAction
	StopAction
	TouchAction
	UnlockAction
	UpdateAction
	UpgradeAction
)

var (
	ErrUnknownAction = errors.New("action unknown")

	ActionNames = map[Action]string{
		AddAction:     "add",
		CreateAction:  "create",
		DisableAction: "disable",
		EnableAction:  "enable",
		InstallAction: "install",
		LockAction:    "lock",
		NilAction:     "nil",
		NothingAction: "nothing",
		ReloadAction:  "reload",
		RestartAction: "restart",
		RemoveAction:  "remove",
		StartAction:   "start",
		StopAction:    "stop",
		RunAction:     "run",
		TouchAction:   "touch",
		UnlockAction:  "unlock",
		UpdateAction:  "update",
		UpgradeAction: "upgrade",
	}

	DelayedNotify taskRunSet = taskRunSet{}
)

type Action int
type ActionMethods map[Action]ActionFunc
type ActionFunc func() (bool, error)
type GuardFunc func() (bool, error)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool

	props  *Properties
	notify actionTaskRunSet
}

type actionTaskRunSet map[Action]map[string]func()
type taskRunSet map[string]func()

func (d *taskRunSet) Run() {
	for _, f := range *d {
		f()
	}
	//clear the list
	d = &taskRunSet{}
}

type Runner interface {
	Run(props *Properties, actions ...Action) bool
}

type Task interface {
	Runner
	fmt.Stringer
	Properties() *Properties
}

func (a Action) name() (string, bool) {
	s, found := ActionNames[a]
	return s, found
}

func (a Action) String() string {
	s, found := a.name()
	if !found {
		s = "UNKNOWN ACTION"
	}
	return s
}

func (r ActionMethods) actionFunc(a Action) (ActionFunc, bool) {
	f, found := r[a]
	return f, found
}

func (b BaseTask) RunActions(task Task, regActions ActionMethods, runActions []Action) bool {
	b.props = task.Properties()

	if len(runActions) == 0 {
		b.logRunStatus(false, false, "error", task, NilAction, time.Now())
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
		f, found := regActions.actionFunc(a)
		if !found {
			b.handleError(b.Errorf(task, a, ErrUnknownAction))
			return hasRun
		}
		hasRun, err := f()
		b.handleError(err)
		b.logRunStatus(hasRun, canRun, reason, task, a, t)
		b.notifyTasks(a)
	}

	return hasRun
}

func (b *BaseTask) NotifyWhen(notify Task, forAction, whenAction Action, props *Properties, delayed bool) {
	if b.notify == nil {
		b.notify = actionTaskRunSet{}
	}
	if b.notify[whenAction] == nil {
		b.notify[whenAction] = map[string]func(){}
	}
	b.notify[whenAction][fmt.Sprintf("%s:%s", notify, forAction)] = func() {
		if delayed {
			DelayedNotify[fmt.Sprintf("%s:%s", notify, forAction)] = func() { notify.Run(props, forAction) }
		} else {
			notify.Run(props, forAction)
		}
	}
}

func (b BaseTask) notifyTasks(action Action) {
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

func (b BaseTask) logRunStatus(hasRun, canRun bool, reason string, task Task, action Action, startTime time.Time) {
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

func (b BaseTask) Errorf(task fmt.Stringer, action Action, err error) error {
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

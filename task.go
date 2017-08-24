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

	DelayedSubscribers delayedSubcribers = delayedSubcribers{}
)

type Action int
type ActionMethods map[Action]ActionFunc
type ActionFunc func() (bool, error)
type GuardFunc func() (bool, error)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	Subscribers map[Action][]func()
	ContOnError bool

	props *Properties
}

type delayedSubcribers []func()

func (d *delayedSubcribers) Run() {
	for _, f := range *d {
		f()
	}
	//clear the list
	d = &delayedSubcribers{}
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
		b.handleError(fmt.Errorf("unable to run %s, no actions specified", task))
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
			b.handleError(fmt.Errorf("%s %s", a, ErrUnknownAction))
			return hasRun
		}
		hasRun, err := f()
		b.handleError(err)
		b.logRunStatus(hasRun, canRun, reason, task, a, t)
		b.notify(a)
	}

	return hasRun
}

func (b *BaseTask) AddSubscriber(task Task, action Action, props *Properties, delayed bool) {
	if b.Subscribers == nil {
		b.Subscribers = map[Action][]func(){}
	}
	b.Subscribers[action] = append(
		b.Subscribers[action],
		func() {
			if delayed {
				DelayedSubscribers = append(DelayedSubscribers, func() { task.Run(props, action) })
			} else {
				task.Run(props, action)
			}
		},
	)
}

func (b BaseTask) notify(action Action) {
	funcs, found := b.Subscribers[action]
	if found {
		for _, f := range funcs {
			f()
		}
	}
}

func (b BaseTask) logRunStatus(hasRun, canRun bool, reason string, t Task, action Action, startTime time.Time) {
	status := ""
	switch {
	case !canRun && reason != "":
		status = fmt.Sprintf("(skipped %s)", reason)
	case !canRun:
		status = fmt.Sprintf("(skipped)", reason)
	case hasRun && reason != "":
		status = fmt.Sprintf("(run %s)", reason)
	case hasRun:
		status = fmt.Sprintf("(run)")
	default:
		status = "(up to date)"
	}

	Log.Printf("  * %s: %s %s %s", t, action, status, time.Since(startTime))
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

func (b BaseTask) handleError(err error) {
	switch {
	case err == nil:
	case !b.ContOnError:
		Log.Panic(err)
	default:
		Log.Printf("    ! %s", err)
	}
}

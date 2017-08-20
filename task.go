package gopack

import (
	"errors"
	"fmt"
	"log"
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

	ErrUnknownAction = errors.New("action unknown")
)

type Action int
type ActionMethods map[Action]ActionFunc
type ActionFunc func() (bool, error)
type GuardFunc func() (bool, error)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	Notify      map[Action][]func()
	ContOnError bool
	logger      *log.Logger
	props       *Properties
}

type Runner interface {
	Run(props *Properties, logger *log.Logger, actions ...Action) bool
}

type Task interface {
	Runner
	fmt.Stringer
	Logger() *log.Logger
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
	b.logger = task.Logger()
	b.props = task.Properties()

	if b.logger == nil {
		b.logger.Panicf("logger is nil for task %s", task)
	}

	if len(runActions) == 0 {
		b.handleError(fmt.Errorf("unable to run %s, no actions specified", task))
		return false
	}

	hasRun := false
	t := time.Now()
	if !b.canRun() {
		b.logRunStatus(hasRun, task, runActions[0], t)
		return hasRun
	}

	for _, a := range runActions {
		if b.runAction(task, regActions, a) {
			hasRun = true
			b.notify(a)
		}
	}

	return hasRun
}

func (b BaseTask) runAction(task Task, regActions ActionMethods, a Action) bool {
	hasRun := false
	t := time.Now()

	f, found := regActions.actionFunc(a)
	if !found {
		b.handleError(fmt.Errorf("%s %s", a, ErrUnknownAction))
		return hasRun
	}

	hasRun, err := f()
	b.handleError(err)
	b.logRunStatus(hasRun, task, a, t)

	return hasRun
}

func (b BaseTask) notify(action Action) {
	funcs, found := b.Notify[action]
	if found {
		for _, f := range funcs {
			f()
		}
	}
}

func (b BaseTask) logRunStatus(hasRun bool, t Task, action Action, startTime time.Time) {
	status := "[NOT RUN]"
	if hasRun {
		status = "[RUN]"
	}
	b.logger.Printf("%s %s %8s %10s\n", t, action, status, time.Since(startTime))
}

func (b BaseTask) canRun() bool {
	var err error
	run := true
	if b.OnlyIf != nil {
		run, err = b.OnlyIf()
		b.handleError(err)
	}
	if b.NotIf != nil {
		run, err = b.NotIf()
		run = !run
		b.handleError(err)
	}
	return run
}

func (b BaseTask) handleError(err error) {
	switch {
	case err == nil:
	case !b.ContOnError:
		b.logger.Panic(err)
	default:
		b.logger.Print(err)
	}
}

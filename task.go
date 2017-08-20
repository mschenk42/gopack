package gopack

import (
	"errors"
	"fmt"
	"log"
	"os"
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
type ActionFunc func(p Properties, logger *log.Logger) (bool, error)
type GuardFunc func(p Properties, logger *log.Logger) (bool, error)

type BaseTask struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	Notify      map[Action][]func()
	ContOnError bool
	Logger      *log.Logger
}

type Runner interface {
	Run(props Properties, logger *log.Logger, actions ...Action) bool
}

type Task interface {
	Runner
	fmt.Stringer
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

func (b BaseTask) RunActions(
	task Task, regActions ActionMethods,
	runActions []Action,
	props Properties,
	logger *log.Logger) bool {

	switch {
	case b.Logger != nil:
	case logger != nil:
		b.Logger = logger
	default:
		b.Logger = log.New(os.Stdout, "", 0)
	}

	if len(runActions) == 0 {
		b.handleError(fmt.Errorf("unable to run %s, no actions specified", task))
		return false
	}

	hasRun := false
	t := time.Now()
	if !b.canRun(props) {
		b.logRunStatus(hasRun, task, runActions[0], t)
		return hasRun
	}

	for _, a := range runActions {
		if b.runAction(task, regActions, a, props) {
			hasRun = true
			b.notify(a)
		}
	}

	return hasRun
}

func (b BaseTask) runAction(
	task Task, regActions ActionMethods,
	a Action,
	props Properties) bool {

	hasRun := false
	t := time.Now()
	f, found := regActions.actionFunc(a)
	if !found {
		b.handleError(fmt.Errorf("%s %s", a, ErrUnknownAction))
		return hasRun
	}

	hasRun, err := f(props, b.Logger)
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
	b.Logger.Printf("%s %s %8s %10s\n", t, action, status, time.Since(startTime))
}

func (b BaseTask) canRun(props Properties) bool {
	var err error
	run := true
	switch {
	case b.OnlyIf != nil:
		run, err = b.OnlyIf(props, b.Logger)
		b.handleError(err)
	case b.NotIf != nil:
		run, err = b.NotIf(props, b.Logger)
		run = !run
		b.handleError(err)
	}
	return run
}

func (b BaseTask) handleError(err error) {
	switch {
	case err == nil:
	case !b.ContOnError:
		b.Logger.Panic(err)
	default:
		b.Logger.Println(err)
	}
}

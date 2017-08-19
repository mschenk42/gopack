package task

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	Add Action = iota
	Create
	Disable
	Enable
	Install
	Lock
	Nothing
	Remove
	Run
	Reload
	Restart
	Start
	Stop
	Touch
	Unlock
	Update
	Upgrade
)

var (
	ActionNames = map[Action]string{
		Add:     "add",
		Create:  "create",
		Disable: "disable",
		Enable:  "enable",
		Install: "install",
		Lock:    "lock",
		Nothing: "nothing",
		Reload:  "reload",
		Restart: "restart",
		Remove:  "remove",
		Start:   "start",
		Stop:    "stop",
		Run:     "run",
		Touch:   "touch",
		Unlock:  "unlock",
		Update:  "update",
		Upgrade: "upgrade",
	}

	ErrUnknownAction = errors.New("action unknown")
)

type Action int
type ActionFunc func(p Properties, logger *log.Logger) (bool, error)
type ActionMethods map[Action]ActionFunc
type GuardFunc func(p Properties, logger *log.Logger) (bool, error)

type Base struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
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

func (b Base) RunActions(
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
		b.notRun(task, runActions[0], t)
		return hasRun
	}

	for _, a := range runActions {
		if b.runAction(task, regActions, a, props) {
			hasRun = true
		}
	}

	return hasRun
}

func (b Base) runAction(
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

	if hasRun {
		b.didRun(task, a, t)
	} else {
		b.notRun(task, a, t)
	}

	return hasRun
}

func (b Base) notRun(t Task, action Action, startTime time.Time) {
	b.Logger.Printf("%s %s %8s %10s\n", t, action, "Not-Run", time.Since(startTime))
}

func (b Base) didRun(t Task, action Action, startTime time.Time) {
	b.Logger.Printf("%s %s %8s %10s\n", t, action, "Did-Run", time.Since(startTime))
}

func (b Base) canRun(props Properties) bool {
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

func (b Base) handleError(err error) {
	switch {
	case err == nil:
		return
	case !b.ContOnError:
		b.Logger.Panic(err)
	default:
		b.Logger.Println(err)
	}
}

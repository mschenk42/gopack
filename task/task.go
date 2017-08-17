package task

import (
	"errors"
	"fmt"
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

	Stdout *os.File = os.Stdout
	Stderr *os.File = os.Stderr
)

type Action int
type ActionFunc func(p Properties) (bool, error)
type ActionMethods map[Action]ActionFunc
type GuardFunc func(p Properties) (bool, error)

type Base struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool
}

type Runner interface {
	Run(props Properties, actions ...Action) bool
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
	props Properties) bool {

	hasRun := false
	t := time.Now()
	if !b.canRun(props) {
		b.notRun(task, Nothing, t)
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
		b.handleError(false, fmt.Errorf("%s %s", a, ErrUnknownAction))
		return hasRun
	}

	hasRun, err := f(props)
	b.handleError(hasRun, err)

	if hasRun {
		b.didRun(task, a, t)
	} else {
		b.notRun(task, a, t)
	}

	return hasRun
}

func (b Base) notRun(t Task, action Action, startTime time.Time) {
	b.info(fmt.Sprintf("%s %s %8s %10s\n", t, action, "Not-Run", time.Since(startTime)))
}

func (b Base) didRun(t Task, action Action, startTime time.Time) {
	b.info(fmt.Sprintf("%s %s %8s %10s\n", t, action, "Did-Run", time.Since(startTime)))
}

func (b Base) canRun(props Properties) bool {
	var err error
	run := true
	switch {
	case b.OnlyIf != nil:
		run, err = b.OnlyIf(props)
		b.handleError(false, err)
	case b.NotIf != nil:
		run, err = b.NotIf(props)
		run = !run
		b.handleError(false, err)
	}
	return run
}

func (b Base) info(s string) {
	fmt.Fprintf(Stdout, s)
}

func (b Base) error(s string) {
	fmt.Fprintf(Stderr, s)
}

func (b Base) handleError(hasRun bool, err error) {
	if err != nil {
		if !b.ContOnError {
			panic(err)
		}
		fmt.Fprintf(Stderr, "%s\n", err)
	}
}

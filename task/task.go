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

type Exec struct {
	OnlyIf      GuardFunc
	NotIf       GuardFunc
	ContOnError bool
	hasRun      bool
	subscribers map[string]notifyAction
}

type notifyAction struct {
	task   Runner
	action Action
	props  Properties
}

type Notifier interface {
	Notify()
}

type Registerer interface {
	Register(r Runner, action ...Action)
}

type Runner interface {
	Notifier
	Run(props Properties, actions ...Action) Runner
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

func (e *Exec) RunActions(
	task Runner,
	regActions ActionMethods,
	runActions []Action,
	props Properties) {

	t := time.Now()
	if !e.canRun(props) {
		e.notRun(task, Nothing, t)
		return
	}

	for _, a := range runActions {
		if e.runAction(task, regActions, a, props) {
			e.hasRun = true
		}
	}
}

func (e *Exec) runAction(
	task Runner,
	regActions ActionMethods,
	a Action,
	props Properties) bool {

	t := time.Now()
	f, found := regActions.actionFunc(a)
	if !found {
		e.handleError(false, fmt.Errorf("%s %s", a, ErrUnknownAction))
		return false
	}

	hasRun, err := f(props)
	e.handleError(hasRun, err)

	if hasRun {
		e.didRun(task, a, t)
	} else {
		e.notRun(task, a, t)
	}

	return hasRun
}

func (e *Exec) DelayNotify(subscriber Runner, a Action, p Properties) {
	if e.subscribers == nil {
		e.subscribers = map[string]notifyAction{}
	}
	e.subscribers[fmt.Sprintf("%s", subscriber)] = notifyAction{task: subscriber, action: a, props: p}
}

func (e *Exec) Notify() {
	if !e.hasRun {
		return
	}
	for _, s := range e.subscribers {
		s.task.Run(s.props, s.action)
	}
}

func (e *Exec) notRun(r Runner, a Action, t time.Time) {
	e.info(fmt.Sprintf("%s %s %8s %10s\n", r, a, "Not-Run", time.Since(t)))
}

func (e *Exec) didRun(r Runner, a Action, t time.Time) {
	e.info(fmt.Sprintf("%s %s %8s %10s\n", r, a, "Did-Run", time.Since(t)))
}

func (e *Exec) canRun(props Properties) bool {
	var err error
	run := true
	switch {
	case e.OnlyIf != nil:
		run, err = e.OnlyIf(props)
		e.handleError(false, err)
	case e.NotIf != nil:
		run, err = e.NotIf(props)
		run = !run
		e.handleError(false, err)
	}
	return run
}

func (e *Exec) info(s string) {
	fmt.Fprintf(Stdout, s)
}

func (e *Exec) error(s string) {
	fmt.Fprintf(Stderr, s)
}

func (e *Exec) handleError(hasRun bool, err error) {
	if err != nil {
		if !e.ContOnError {
			panic(err)
		}
		fmt.Fprintf(Stderr, "%s\n", err)
	}
}

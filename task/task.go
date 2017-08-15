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
	hasRun      bool
	subscribers map[string]notifyAction
}

type notifyAction struct {
	task   Task
	action Action
	props  Properties
}

type Identifier interface {
	// should be a unique name for the task
	ID() string
}

type Notifier interface {
	Notify()
}

type Registerer interface {
	Register(r Task, action ...Action)
}

type Runner interface {
	Run(props Properties, actions ...Action) Task
}

type Task interface {
	Identifier
	Notifier
	Runner
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

func (b *Base) RunActions(
	task Task,
	regActions ActionMethods,
	runActions []Action,
	props Properties) {

	t := time.Now()
	if !b.canRun(props) {
		b.notRun(task, Nothing, t)
		return
	}

	for _, a := range runActions {
		b.runAction(task, regActions, a, props)
	}
}

func (b *Base) runAction(
	task Task,
	regActions ActionMethods,
	a Action,
	props Properties) {

	t := time.Now()
	f, found := regActions.actionFunc(a)
	if !found {
		b.handleError(false, fmt.Errorf("%s %s", a, ErrUnknownAction))
		return
	}

	var err error
	b.hasRun, err = f(props)
	b.handleError(b.hasRun, err)

	if b.hasRun {
		b.didRun(task, a, t)
	} else {
		b.notRun(task, a, t)
	}

}

func (b *Base) DelayNotify(subscriber Task, action Action, props Properties) {
	if b.subscribers == nil {
		b.subscribers = map[string]notifyAction{}
	}
	b.subscribers[subscriber.ID()] = notifyAction{task: subscriber, action: action, props: props}
}

func (b *Base) Notify() {
	if !b.hasRun {
		return
	}
	for _, s := range b.subscribers {
		s.task.Run(s.props, s.action)
	}
}

func (b *Base) notRun(t Task, action Action, startTime time.Time) {
	b.info(fmt.Sprintf("%s %s %8s %10s\n", t, action, "Not-Run", time.Since(startTime)))
}

func (b *Base) didRun(t Task, action Action, startTime time.Time) {
	b.info(fmt.Sprintf("%s %s %8s %10s\n", t, action, "Did-Run", time.Since(startTime)))
}

func (b *Base) canRun(props Properties) bool {
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

func (b *Base) info(s string) {
	fmt.Fprintf(Stdout, s)
}

func (b *Base) error(s string) {
	fmt.Fprintf(Stderr, s)
}

func (b *Base) handleError(hasRun bool, err error) {
	if err != nil {
		if !b.ContOnError {
			panic(err)
		}
		fmt.Fprintf(Stderr, "%s\n", err)
	}
}

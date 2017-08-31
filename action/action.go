package action

import "errors"

const (
	Add Enum = iota
	Create
	Disable
	Enable
	Install
	Lock
	Nil
	Nothing
	Reload
	Remove
	Restart
	Run
	Start
	Stop
	Touch
	Unlock
	Update
	Upgrade
)

var (
	ErrActionNotRegistered = errors.New("action not registered with task")

	Names = map[Enum]string{
		Add:     "add",
		Create:  "create",
		Disable: "disable",
		Enable:  "enable",
		Install: "install",
		Lock:    "lock",
		Nil:     "nil",
		Nothing: "nothing",
		Reload:  "reload",
		Remove:  "remove",
		Restart: "restart",
		Run:     "run",
		Start:   "start",
		Stop:    "stop",
		Touch:   "touch",
		Unlock:  "unlock",
		Update:  "update",
		Upgrade: "upgrade",
	}
)

type Enum int
type Methods map[Enum]methodFunc
type methodFunc func() (bool, error)

func (a Enum) name() (string, bool) {
	s, found := Names[a]
	return s, found
}

func (a Enum) String() string {
	s, found := a.name()
	if !found {
		s = "UNKNOWN ACTION"
	}
	return s
}

func (m Methods) Method(a Enum) (methodFunc, bool) {
	f, found := m[a]
	return f, found
}

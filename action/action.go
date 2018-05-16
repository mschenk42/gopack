package action

import "fmt"

const (
	Add Name = iota
	Copy
	Create
	Disable
	Down
	Download
	Enable
	Install
	Lock
	Nil
	Move
	Nothing
	Reload
	Remove
	Restart
	Run
	Start
	Stop
	Touch
	Unlock
	Up
	Update
	Upgrade
)

var (
	names = map[Name]string{
		Add:      "add",
		Copy:     "copy",
		Create:   "create",
		Disable:  "disable",
		Down:     "down",
		Download: "download",
		Enable:   "enable",
		Install:  "install",
		Lock:     "lock",
		Move:     "move",
		Nil:      "nil",
		Nothing:  "nothing",
		Reload:   "reload",
		Remove:   "remove",
		Restart:  "restart",
		Run:      "run",
		Start:    "start",
		Stop:     "stop",
		Touch:    "touch",
		Unlock:   "unlock",
		Up:       "up",
		Update:   "update",
		Upgrade:  "upgrade",
	}
)

type Name int
type Func func() (bool, error)
type Funcs map[Name]Func

func (a Name) name() (string, bool) {
	x, found := names[a]
	return x, found
}

func (a Name) String() string {
	x, found := a.name()
	if !found {
		x = fmt.Sprintf("Unknown action %d", a)
	}
	return x
}

func (m Funcs) Func(a Name) (Func, bool) {
	x, found := m[a]
	return x, found
}

// NewSlice is a helper method for creating action slices
func NewSlice(a ...Name) []Name {
	return a
}

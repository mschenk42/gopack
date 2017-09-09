package action

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
	names = map[Enum]string{
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
type Methods map[Enum]MethodFunc
type MethodFunc func() (bool, error)

func (a Enum) name() (string, bool) {
	x, found := names[a]
	return x, found
}

func (a Enum) String() string {
	x, found := a.name()
	if !found {
		x = "UNKNOWN ACTION"
	}
	return x
}

func (m Methods) Method(a Enum) (MethodFunc, bool) {
	x, found := m[a]
	return x, found
}

// NewSlice is a helper method for creating action slices
func NewSlice(a ...Enum) []Enum {
	return a
}

package task

import (
	"fmt"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

// Func ...
type Func struct {
	ActionFunc action.Func

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (f Func) Run(runActions ...action.Name) gopack.ActionRunStatus {
	f.setDefaults()
	return f.RunActions(&f, f.registerActions(), runActions)
}

func (f Func) registerActions() action.Funcs {
	return action.Funcs{
		action.Run: f.run,
	}
}

func (f *Func) setDefaults() {
}

// String returns a string which identifies the task with it's property values
func (f Func) String() string {
	return fmt.Sprintf("func %p", f.ActionFunc)
}

func (f Func) run() (bool, error) {
	return f.ActionFunc()
}

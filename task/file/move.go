package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/mschenk42/gopack/task"
)

// Move ...
type Move struct {
	From  string
	To    string
	Owner string
	Group string
	Perm  os.FileMode

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (m Move) Run(runActions ...action.Name) gopack.ActionRunStatus {
	m.setDefaults()
	return m.RunActions(&m, m.registerActions(), runActions)
}

func (m Move) registerActions() action.Funcs {
	return action.Funcs{
		action.Run: m.run,
	}
}

func (m *Move) setDefaults() {
}

// String returns a string which identifies the task with it's property values
func (m Move) String() string {
	return fmt.Sprintf("move %s %s %s %s %s", m.From, m.To, m.Owner, m.Group, m.Perm)
}

func (m Move) run() (bool, error) {
	_, exists, err := task.Fexists(m.From)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	b, err := ioutil.ReadFile(m.From)
	if err != nil {
		return false, err
	}
	if err := ioutil.WriteFile(m.To, b, m.Perm); err != nil {
		return false, err
	}
	if _, err := task.Chown(m.To, m.Owner, m.Group); err != nil {
		return false, err
	}
	if err := os.Remove(m.From); err != nil {
		return false, err
	}
	return true, nil
}

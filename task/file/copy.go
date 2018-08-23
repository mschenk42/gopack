package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/mschenk42/gopack/task"
)

// Copy ...
type Copy struct {
	From  string
	To    string
	Owner string
	Group string
	Perm  os.FileMode

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (c Copy) Run(runActions ...action.Name) gopack.ActionRunStatus {
	c.setDefaults()
	return c.RunActions(&c, c.registerActions(), runActions)
}

func (c Copy) registerActions() action.Funcs {
	return action.Funcs{
		action.Run: c.run,
	}
}

func (c *Copy) setDefaults() {
}

// String returns a string which identifies the task with it's property values
func (c Copy) String() string {
	return fmt.Sprintf("copy %s %s %s %s %s", c.From, c.To, c.Owner, c.Group, c.Perm)
}

func (c Copy) run() (bool, error) {
	_, exists, err := task.Fexists(c.From)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	b, err := ioutil.ReadFile(c.From)
	if err != nil {
		return false, err
	}
	if err := ioutil.WriteFile(c.To, b, c.Perm); err != nil {
		return false, err
	}
	if _, err := task.Chown(c.To, c.Owner, c.Group); err != nil {
		return false, err
	}
	return true, nil
}

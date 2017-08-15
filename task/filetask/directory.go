package filetask

import (
	"fmt"
	"os"

	"github.com/mschenk42/mincfg/task"
)

type Directory struct {
	Name  string
	Path  string
	Group string
	Owner string
	Perm  os.FileMode
	task.Exec
}

func (d Directory) Register(x task.Registerer, runActions ...task.Action) Directory {
	x.Register(&d, runActions...)
	return d
}

func (d Directory) Run(props task.Properties, runActions ...task.Action) task.Task {
	regActions := task.ActionMethods{
		task.Create: d.create,
		task.Remove: d.remove,
	}
	d.Exec.RunActions(&d, regActions, runActions, props)
	return &d
}

func (d Directory) ID() string {
	return d.Name
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s", d.Path)
}

func (d Directory) create(props task.Properties) (bool, error) {
	x, err := d.exists()
	switch {
	case err != nil:
		return false, err
	case x:
		return false, nil
	default:
		err := os.MkdirAll(d.Path, d.Perm)
		return true, err
	}
}

func (d Directory) remove(props task.Properties) (bool, error) {
	x, err := d.exists()
	switch {
	case err != nil:
		return false, err
	case !x:
		return false, nil
	default:
		err := os.Remove(d.Path)
		return true, err
	}
}

func (d Directory) exists() (bool, error) {
	_, err := os.Stat(d.Path)
	switch {
	case err != nil:
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	default:
		return true, nil
	}
}

package task

import (
	"fmt"
	"os"

	"github.com/mschenk42/gopack"
)

type Directory struct {
	Path  string
	Owner string
	Group string
	Mode  os.FileMode

	gopack.BaseTask
}

func (d Directory) Run(runActions ...gopack.Action) bool {
	d.setDefaults()
	return d.RunActions(&d, d.registerActions(), runActions)
}

func (d Directory) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: d.create,
		gopack.RemoveAction: d.remove,
	}
}

func (d *Directory) setDefaults() {
	if d.Mode == 0 {
		d.Mode = 0755
	}
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s %s %s %s", d.Path, d.Owner, d.Group, d.Mode)
}

func (d Directory) create() (bool, error) {
	var (
		err          error
		found        bool
		chgOwnership bool
		chgDirectory bool
		chgMode      bool
		fi           os.FileInfo
	)

	if fi, found, err = fexists(d.Path); err != nil {
		return false, d.TaskError(d, gopack.CreateAction, err)
	}
	if !found {
		chgDirectory = true
		if err = os.MkdirAll(d.Path, d.Mode); err != nil {
			return false, d.TaskError(d, gopack.CreateAction, err)
		}
	} else {
		if fi.Mode().Perm() != d.Mode.Perm() {
			os.Chmod(d.Path, d.Mode)
			chgMode = true
		}
	}

	if d.Owner == "" && d.Group == "" {
		return chgDirectory || chgOwnership || chgMode, nil
	}
	if chgOwnership, err = chown(d.Path, d.Owner, d.Group); err != nil {
		return false, d.TaskError(d, gopack.CreateAction, err)
	}
	return chgDirectory || chgOwnership || chgMode, nil
}

func (d Directory) remove() (bool, error) {
	var (
		found bool
		err   error
	)
	if _, found, err = fexists(d.Path); err != nil {
		return false, d.TaskError(d, gopack.CreateAction, err)
	}
	if !found {
		return false, nil
	}
	//TODO: optionally allow RemoveAll
	err = os.Remove(d.Path)
	return true, d.TaskError(d, gopack.RemoveAction, err)
}

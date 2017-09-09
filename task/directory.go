package task

import (
	"fmt"
	"os"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

type Directory struct {
	Path  string
	Owner string
	Group string
	Mode  os.FileMode

	gopack.BaseTask
}

func (d Directory) Run(runActions ...action.Enum) gopack.ActionRunStatus {
	d.setDefaults()
	return d.RunActions(&d, d.registerActions(), runActions)
}

func (d Directory) registerActions() action.Methods {
	return action.Methods{
		action.Create: d.create,
		action.Remove: d.remove,
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

	if fi, found, err = Fexists(d.Path); err != nil {
		return false, err
	}
	if !found {
		chgDirectory = true
		if err = os.MkdirAll(d.Path, d.Mode); err != nil {
			return false, err
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
	if chgOwnership, err = Chown(d.Path, d.Owner, d.Group); err != nil {
		return false, err
	}
	return chgDirectory || chgOwnership || chgMode, nil
}

func (d Directory) remove() (bool, error) {
	var (
		found bool
		err   error
	)
	if _, found, err = Fexists(d.Path); err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	//TODO: optionally allow RemoveAll
	err = os.Remove(d.Path)
	return true, err
}

package task

import (
	"fmt"
	"log"
	"os"

	"github.com/mschenk42/gopack"
)

type Directory struct {
	Name  string
	Path  string
	Group string
	Owner string
	Perm  os.FileMode

	logger     *log.Logger
	properties *gopack.Properties
	defaults   *gopack.Properties
	gopack.BaseTask
}

func (d Directory) Run(props *gopack.Properties, logger *log.Logger, runActions ...gopack.Action) bool {
	d.logger = logger
	d.properties = props
	d.setDefaults()
	return d.BaseTask.RunActions(&d, d.registerActions(), runActions)
}

func (d Directory) Logger() *log.Logger {
	return d.logger
}

func (d Directory) Properties() *gopack.Properties {
	return d.properties
}

func (d Directory) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: d.create,
		gopack.RemoveAction: d.remove,
	}
}

func (d *Directory) setDefaults() {
	switch {
	case d.Perm == 0:
		d.Perm = 0755
	}
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s", d.Path)
}

func (d Directory) create() (bool, error) {
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

func (d Directory) remove() (bool, error) {
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

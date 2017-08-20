package task

import (
	"fmt"
	"log"
	"os"

	"github.com/mschenk42/gopack"
)

type Directory struct {
	Name     string
	Path     string
	Group    string
	Owner    string
	Perm     os.FileMode
	defaults gopack.Properties
	gopack.BaseTask
}

func (d Directory) Run(props gopack.Properties, logger *log.Logger, runActions ...gopack.Action) bool {
	regActions := gopack.ActionMethods{
		gopack.CreateAction: d.create,
		gopack.RemoveAction: d.remove,
	}

	d.defaults = gopack.Properties{
		"perm": 0755,
	}

	return d.BaseTask.RunActions(&d, regActions, runActions, props, logger)
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s", d.Path)
}

func (d Directory) create(props gopack.Properties, logger *log.Logger) (bool, error) {
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

func (d Directory) remove(props gopack.Properties, logger *log.Logger) (bool, error) {
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

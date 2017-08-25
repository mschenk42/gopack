package task

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/mschenk42/gopack"
)

type Directory struct {
	Path  string
	Group string
	Owner string
	Perm  os.FileMode

	props *gopack.Properties
	gopack.BaseTask
}

func (d Directory) Run(props *gopack.Properties, runActions ...gopack.Action) bool {
	d.props = props
	d.setDefaults()
	return d.BaseTask.RunActions(&d, d.registerActions(), runActions)
}

func (d Directory) Properties() *gopack.Properties {
	return d.props
}

func (d Directory) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: d.create,
		gopack.RemoveAction: d.remove,
	}
}

func (d *Directory) setDefaults() {
	if d.Perm == 0 {
		d.Perm = 0755
	}
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s Owner:%s Group:%s Perm:%s", d.Path, d.Owner, d.Group, d.Perm)
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
		if err != nil {
			return false, err
		}

		if d.Owner == "" && d.Group == "" {
			return true, err
		}

		var u *user.User
		var gid, uid int

		if d.Owner == "" {
			u, err = user.Current()
			if err != nil {
				return false, err
			}
		} else {
			u, err = user.Lookup(d.Owner)
			if err != nil {
				return false, err
			}
			uid, err = strconv.Atoi(u.Uid)
			if err != nil {
				return false, err
			}
		}

		if d.Group == "" {
			gid, err = strconv.Atoi(u.Gid)
			if err != nil {
				return false, err
			}
		} else {
			g, err := user.LookupGroup(d.Group)
			if err != nil {
				return false, err
			}
			gid, err = strconv.Atoi(g.Gid)
			if err != nil {
				return false, err
			}
		}

		err = os.Chown(d.Path, uid, gid)
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
		//TODO: optionally allow RemoveAll
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

package task

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/mschenk42/gopack"
)

type Directory struct {
	Path  string
	Owner string
	Group string
	Mode  os.FileMode

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
	)

	if found, err = fexists(d.Path); err != nil {
		return false, d.Errorf(d, gopack.CreateAction, err)
	}
	if !found {
		chgDirectory = true
		if err = os.MkdirAll(d.Path, d.Mode); err != nil {
			return false, d.Errorf(d, gopack.CreateAction, err)
		}
	}

	if d.Owner == "" && d.Group == "" {
		return chgDirectory, nil
	}
	if chgOwnership, err = chown(d.Path, d.Owner, d.Group); err != nil {
		return false, d.Errorf(d, gopack.CreateAction, err)
	}
	return chgDirectory || chgOwnership, nil
}

func (d Directory) remove() (bool, error) {
	var (
		found bool
		err   error
	)
	if found, err = fexists(d.Path); err != nil {
		return false, d.Errorf(d, gopack.CreateAction, err)
	}
	if !found {
		return false, nil
	}
	//TODO: optionally allow RemoveAll
	err = os.Remove(d.Path)
	return true, d.Errorf(d, gopack.RemoveAction, err)
}

func chown(path, owner, group string) (bool, error) {
	var (
		err      error
		u        *user.User
		g        *user.Group
		gid, uid int
	)

	// use current user if no owner provided
	if owner == "" {
		if u, err = user.Current(); err != nil {
			return false, err
		}
	} else {
		if u, err = user.Lookup(owner); err != nil {
			return false, err
		}
	}
	if uid, err = strconv.Atoi(u.Uid); err != nil {
		return false, err
	}

	// use user's group if no group provided
	if group == "" {
		if gid, err = strconv.Atoi(u.Gid); err != nil {
			return false, err
		}
	} else {
		if g, err = user.LookupGroup(group); err != nil {
			return false, err
		}
		if gid, err = strconv.Atoi(g.Gid); err != nil {
			return false, err
		}
	}

	// check if ownership is differrent then provided
	var uidNow, gidNow int
	if uidNow, gidNow, err = fownership(path); err != nil {
		return false, err
	}
	if uid == uidNow && gid == gidNow {
		return false, nil
	}

	// set ownership
	if err = os.Chown(path, uid, gid); err != nil {
		return false, err
	}

	return true, nil
}

func fexists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func fownership(path string) (int, int, error) {
	var (
		fi  os.FileInfo
		err error
	)
	if fi, err = os.Stat(path); err != nil {
		return 0, 0, err
	}
	uid := fi.Sys().(*syscall.Stat_t).Uid
	gid := fi.Sys().(*syscall.Stat_t).Gid
	return int(uid), int(gid), nil
}

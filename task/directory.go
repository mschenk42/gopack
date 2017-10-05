package task

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

type Directory struct {
	Path  string
	Owner string
	Group string
	Perm  os.FileMode

	gopack.BaseTask
}

func (d Directory) Run(runActions ...action.Name) gopack.ActionRunStatus {
	d.setDefaults()
	return d.RunActions(&d, d.registerActions(), runActions)
}

func (d Directory) registerActions() action.Funcs {
	return action.Funcs{
		action.Create: d.create,
		action.Remove: d.remove,
	}
}

func (d *Directory) setDefaults() {
	if d.Perm == 0 {
		d.Perm = 0755
	}
}

func (d Directory) String() string {
	return fmt.Sprintf("directory %s %s %s %s", d.Path, d.Owner, d.Group, d.Perm)
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
		if err = os.MkdirAll(d.Path, d.Perm); err != nil {
			return false, err
		}
	} else {
		if fi.Mode().Perm() != d.Perm.Perm() {
			os.Chmod(d.Path, d.Perm)
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

func Chown(path, owner, group string) (bool, error) {
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
	var (
		fi     os.FileInfo
		uidNow int
		gidNow int
	)
	if fi, err = os.Stat(path); err != nil {
		return false, err
	}
	if fi.Sys() != nil {
		uidNow = int(fi.Sys().(*syscall.Stat_t).Uid)
		gidNow = int(fi.Sys().(*syscall.Stat_t).Gid)
	} else {
		return false, fmt.Errorf("syscall is nil for %s", path)
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

package task

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func fexists(path string) (os.FileInfo, bool, error) {
	var (
		err error
		fi  os.FileInfo
	)
	if fi, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fi, false, nil
		}
		return fi, false, err
	}
	return fi, true, nil
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

func execCmd(timeout time.Duration, command string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return b, fmt.Errorf("%s %s", strings.Replace(string(b), "\n", " ", -1), err)
	}
	if ctx.Err() == context.DeadlineExceeded {
		return b, ctx.Err()
	}

	return b, nil
}

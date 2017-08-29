package task

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

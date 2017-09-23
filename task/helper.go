package task

import (
	"os"
)

func Fexists(path string) (os.FileInfo, bool, error) {
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

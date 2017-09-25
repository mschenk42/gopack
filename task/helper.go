package task

import (
	"os"
	"strings"
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

func Redact(strs ...string) string {
	masked := ""
	for _, s := range strs {
		masked += " " + strings.Repeat("*", len(s))
	}
	return masked
}

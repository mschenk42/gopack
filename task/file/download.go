package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/mschenk42/gopack/task"
)

// Download ...
type Download struct {
	URL   string
	Path  string
	Owner string
	Group string
	Perm  os.FileMode

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func (d Download) Run(runActions ...action.Name) gopack.ActionRunStatus {
	d.setDefaults()
	return d.RunActions(&d, d.registerActions(), runActions)
}

func (d Download) registerActions() action.Funcs {
	return action.Funcs{
		action.Run: d.create,
	}
}

func (d *Download) setDefaults() {
}

// String returns a string which identifies the task with it's property values
func (d Download) String() string {
	return fmt.Sprintf("download %s %s %s %s %s", d.URL, d.Path, d.Owner, d.Group, d.Perm)
}

func (d Download) create() (bool, error) {
	out, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, d.Perm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	resp, err := http.Get(d.URL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	if _, err := task.Chown(d.Path, d.Owner, d.Group); err != nil {
		return false, err
	}
	return true, nil
}

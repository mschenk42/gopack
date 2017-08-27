package task

import (
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/task"
)

func TestCreateTemplate(t *testing.T) {
	task.Template{
		Name:   "test-template",
		Path:   "/tmp/test/mypack.conf",
		Mode:   0755,
		Source: `log_dir:{{ index . "nginx.log_dir"}}`,
	}.Run(
		pack.Props,
		gopack.CreateAction,
	)
}

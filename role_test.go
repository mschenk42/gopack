package mincfg

import (
	"testing"

	"github.com/mschenk42/mincfg/task"
	"github.com/mschenk42/mincfg/task/filetask"
)

func TestCreateRole(t *testing.T) {
	nginx := &Role{
		Name: "nginx",
		Props: task.Properties{
			"nginx.log.dir": "/tmp/nginx",
		},
	}

	filetask.Directory{
		Path: nginx.Props.Str("nginx.log.dir"),
		Perm: 0755,
	}.Register(
		nginx,
		task.Create,
		task.Remove,
	)

	nginx.Run(nil)
}

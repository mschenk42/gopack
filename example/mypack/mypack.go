package mypack

import (
	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/task"
)

func Run(props *gopack.Properties) {
	pack := gopack.Pack{
		Name: "MyPack",
		Props: &gopack.Properties{
			"nginx.log_dir":   "/var/log/nginx",
			"nginx.cache_dir": "/var/cache",
		},
		RunFunc: run,
	}
	pack.Run(props)
}

func run(pack *gopack.Pack) {

	task.User{
		UserName: "mypack",
		Group:    "sudo",
	}.Run(
		pack.Props,
		gopack.CreateAction,
	)

	task.Directory{
		Path:  "/tmp/test",
		Owner: "mypack",
		Mode:  0755,
	}.Run(
		pack.Props,
		gopack.CreateAction,
	)

	task.Template{
		Name:   "mypack",
		Path:   "/tmp/test/mypack.conf",
		Owner:  "mypack",
		Mode:   0755,
		Source: `log_dir:{{ index . "nginx.log_dir"}}`,
	}.Run(
		pack.Props,
		gopack.CreateAction,
	)

	task.User{
		UserName: "mypack",
	}.Run(
		pack.Props,
		gopack.RemoveAction,
	)
}

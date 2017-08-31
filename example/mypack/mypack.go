package mypack

import (
	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/task"
)

func Run(props *gopack.Properties) {
	pack := gopack.Pack{
		Name: "MyPack",
		Props: &gopack.Properties{
			"mypack.user":   "mypack",
			"mypack.group":  "mypack",
			"nginx.log_dir": "/var/log/nginx",
		},
		RunFunc: run,
	}
	pack.Run(props)
}

func run(pack *gopack.Pack) {

	owner, _ := (*pack.Props).Str("mypack.user")
	group, _ := (*pack.Props).Str("mypack.group")

	task.Group{Name: "mypack"}.Run(gopack.CreateAction)
	task.User{Name: "mypack", Group: "mypack"}.Run(gopack.CreateAction)

	task.Directory{
		Path:  "/tmp/test",
		Owner: owner,
		Group: group,
		Mode:  0755,
	}.Run(gopack.CreateAction)

	task.Template{
		Name:   "mypack",
		Path:   "/tmp/test/mypack.conf",
		Owner:  owner,
		Mode:   0755,
		Source: `log_dir:{{ index . "nginx.log_dir"}}`,
		Props:  pack.Props,
	}.Run(gopack.CreateAction)
}

package mincfg

import (
	"log"
	"os"

	"github.com/mschenk42/mincfg/task"
)

type RunBook struct {
	Name   string
	Props  task.Properties
	Logger *log.Logger
	roles  []Role
}

func (r *RunBook) Register(role Role) {
	if r.roles == nil {
		r.roles = []Role{}
	}
	r.roles = append(r.roles, role)
}

func (r *RunBook) Run(props task.Properties) {
	if r.Logger == nil {
		r.Logger = log.New(os.Stdout, "", 0)
	}

	r.Props.Merge(props)
	for _, role := range r.roles {
		role.Run(r.Props)
	}
}

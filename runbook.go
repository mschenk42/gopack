package mincfg

import "github.com/mschenk42/mincfg/task"

type RunBook struct {
	Name  string
	Props task.Properties
	roles []Role
}

func (r *RunBook) Register(role Role) {
	if r.roles == nil {
		r.roles = []Role{}
	}
	r.roles = append(r.roles, role)
}

func (r *RunBook) Run(props task.Properties) {
	r.Props.Merge(props)
	for _, role := range r.roles {
		role.Run(r.Props)
	}
	r.Notify()
}

func (r *RunBook) Notify() {
	for _, role := range r.roles {
		role.Notify()
	}
}

func (r *RunBook) handleError(err error) {
	if err != nil {
		panic(err)
	}
}

package mincfg

import "github.com/mschenk42/mincfg/task"

type RunBook struct {
	Name  string
	Props task.Properties
	roles []Role
}

func (r *RunBook) Register(x Role) {
	if r.roles == nil {
		r.roles = []Role{}
	}
	r.roles = append(r.roles, x)
}

func (r *RunBook) Run(p task.Properties) {
	r.Props.Merge(p)
	for _, x := range r.roles {
		x.Run(r.Props)
	}
	r.Notify()
}

func (r *RunBook) Notify() {
	for _, x := range r.roles {
		x.Notify()
	}
}

func (r *RunBook) handleError(err error) {
	if err != nil {
		panic(err)
	}
}

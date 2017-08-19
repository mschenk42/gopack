package gopack

import (
	"log"
	"os"
)

type Pack struct {
	Name   string
	Props  Properties
	Logger *log.Logger
	roles  []Role
}

func (p *Pack) Register(role Role) {
	if p.roles == nil {
		p.roles = []Role{}
	}
	p.roles = append(p.roles, role)
}

func (p *Pack) Run(props Properties) {
	if p.Logger == nil {
		p.Logger = log.New(os.Stdout, "", 0)
	}

	p.Props.Merge(props)
	for _, role := range p.roles {
		role.Run(p.Props)
	}
}

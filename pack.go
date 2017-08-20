package gopack

import (
	"log"
	"os"
)

type Pack struct {
	Name    string
	Props   Properties
	Logger  *log.Logger
	RunFunc func()
}

func (p *Pack) Run(props Properties) {
	if p.Logger == nil {
		p.Logger = log.New(os.Stdout, "", 0)
	}
	p.Props.Merge(props)
	p.RunFunc()
}

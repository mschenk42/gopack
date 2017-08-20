package gopack

import (
	"log"
	"os"
)

type Pack struct {
	Name    string
	Props   *Properties
	Logger  *log.Logger
	RunFunc func(pack *Pack)
}

func (p *Pack) Run(props *Properties) {
	if p.Logger == nil {
		p.Logger = log.New(os.Stdout, "", 0)
	}
	if p.RunFunc == nil {
		p.Logger.Panic("Run function nil for pack %s", p.Name)
	}
	p.Props.Merge(props)
	p.RunFunc(p)
}

package gopack

import (
	"log"
	"os"
)

var Logger *log.Logger = log.New(os.Stdout, "", 0)

type Pack struct {
	Name    string
	Props   *Properties
	RunFunc func(pack *Pack)
}

func (p *Pack) Run(props *Properties) {
	if p.RunFunc == nil {
		Logger.Panic("run function nil for pack %s", p.Name)
	}
	p.Props.Merge(props)
	p.RunFunc(p)
}

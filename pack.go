package gopack

import (
	"log"
	"os"
)

var Log Logger = log.New(os.Stdout, "", log.Lmicroseconds)

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Pack struct {
	Name    string
	Props   *Properties
	RunFunc func(pack *Pack)
}

func (p *Pack) Run(props *Properties) {
	if p.RunFunc == nil {
		Log.Panic("run function nil for pack %s", p.Name)
	}
	p.Props.Merge(props)
	p.RunFunc(p)
}

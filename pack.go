package gopack

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Log Logger = log.New(os.Stdout, "", 0)

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
	Name         string
	Props        *Properties
	Redact       []string
	RunFunc      func(pack *Pack)
	NoRunDelayed bool
}

func (p Pack) String() string {
	return fmt.Sprintf("%s", p.Name)
}

func (p *Pack) Run(props *Properties) {
	t := time.Now()
	if p.RunFunc == nil {
		Log.Fatalf("  ! run function nil for pack %s", p.Name)
	}
	p.Props.Merge(props)
	Log.Printf("Pack: %s (start)", p)
	Log.Printf("%s", p.Props.Redact(p.Redact))
	Log.Println("\n  [run]")
	p.RunFunc(p)
	if !p.NoRunDelayed {
		Log.Println("\n  [delayed run]")
		DelayedNotify.Run()
	}
	Log.Printf("\nPack: %s (end) %s\n", p, time.Since(t))
}

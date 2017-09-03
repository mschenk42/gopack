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
	Actions      []string
	NoRunDelayed bool
}

var (
	PackHeaderFormat   string = colorize.Blue("\nPack: %s (%s) %s")
	PackSectionFormat  string = colorize.Blue("\n[%s %s]\n")
	PackErrorFormat    string = colorize.Red("! %s\n")
	PackPropertyFormat string = colorize.Magenta("%s")
)

func (p Pack) String() string {
	return fmt.Sprintf("%s", p.Name)
}

func (p *Pack) Run(props *Properties) {
	t := time.Now()
	if p.RunFunc == nil {
		Log.Fatalf(PackErrorFormat, fmt.Sprintf("run function nil for pack %s", p.Name))
	}
	p.Props.Merge(props)
	Log.Printf(PackHeaderFormat, p, "start", "")
	Log.Printf(PackPropertyFormat, p.Props.Redact(p.Redact))
	Log.Printf(PackSectionFormat, "run tasks for", p)
	p.RunFunc(p)
	if !p.NoRunDelayed {
		Log.Printf(PackSectionFormat, "run delayed tasks for", p)
		DelayedNotify.Run()
	}
	Log.Printf(PackHeaderFormat, p, "end", time.Since(t))
	Log.Print("")
}

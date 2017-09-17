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
	Actions      []string
	ActionMap    map[string]func(p *Pack)
	NoRunDelayed bool
}

var (
	PackHeaderFormat       string = colorize.Blue("\nPack: %s (%s) %s")
	PackSectionFormat      string = colorize.Blue("\n[%s %s]\n")
	PackErrorFormat        string = colorize.Red("! %s\n")
	PackPropertyFormat     string = colorize.Magenta("%s")
	PackActionHeaderFormat string = colorize.Blue("Pack: %s %s (%s) %s")
)

func (p Pack) String() string {
	return fmt.Sprintf("%s", p.Name)
}

func (p *Pack) Run(props *Properties) {
	t := time.Now()
	p.Props.Merge(props)
	Log.Printf(PackHeaderFormat, p, "start", "")
	Log.Printf(PackPropertyFormat, p.Props.Redact(p.Redact))
	Log.Printf(PackSectionFormat, "run actions for", p)

	p.run()
	if !p.NoRunDelayed {
		Log.Printf(PackSectionFormat, "run delayed tasks for", p)
		delayedNotify.Run()
	}

	Log.Printf(PackSectionFormat, "summary of tasks run for", p)
	for _, x := range tasksRun {
		Log.Print(x)
	}

	Log.Printf(PackHeaderFormat, p, "end", time.Since(t))
	Log.Print("")
}

func (p *Pack) run() {
	if len(p.Actions) == 0 {
		p.ActionMap["default"](p)
	}
	for _, action := range p.Actions {
		if f, found := p.ActionMap[action]; found {
			t := time.Now()
			Log.Printf(PackActionHeaderFormat, p, action, "start", "")
			f(p)
			Log.Printf(PackActionHeaderFormat, p, action, "end", time.Since(t))
		} else {
			Log.Fatalf(PackErrorFormat, fmt.Sprintf("pack action %s not found", action))
		}
	}
}

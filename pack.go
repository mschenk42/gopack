package gopack

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mschenk42/gopack/color"
)

var Log = log.New(os.Stdout, "", 0)

type Pack struct {
	Name         string
	Props        *Properties
	Redact       []string
	Actions      []string
	ActionMap    map[string]func(p *Pack)
	NoRunDelayed bool
}

var (
	packHeaderFormat       = color.Blue("\nPack: %s (%s) %s")
	packSectionFormat      = color.Blue("\n[%s %s]\n")
	packErrorFormat        = color.Red("! %s\n")
	packPropertyFormat     = color.Magenta("%s")
	packActionHeaderFormat = color.Blue("Pack: %s %s (%s) %s")
)

func (p Pack) String() string {
	return fmt.Sprintf("%s", p.Name)
}

func (p *Pack) Run(props *Properties) {
	t := time.Now()
	p.Props.Merge(props)
	Log.Printf(packHeaderFormat, p, "start", "")
	Log.Printf(packPropertyFormat, p.Props.Redact(p.Redact))
	Log.Printf(packSectionFormat, "run actions for", p)

	p.run()
	if !p.NoRunDelayed {
		Log.Printf(packSectionFormat, "run delayed tasks for", p)
		delayedNotify.Run()
	}

	Log.Printf(packSectionFormat, "summary of tasks run for", p)
	for _, x := range tasksRun {
		Log.Print(x)
	}

	Log.Printf(packHeaderFormat, p, "end", time.Since(t))
	Log.Print("")
}

func (p *Pack) run() {
	if len(p.Actions) == 0 {
		p.ActionMap["default"](p)
	}
	for _, action := range p.Actions {
		if f, found := p.ActionMap[action]; found {
			t := time.Now()
			Log.Printf(packActionHeaderFormat, p, action, "start", "")
			f(p)
			Log.Printf(packActionHeaderFormat, p, action, "end", time.Since(t))
		} else {
			Log.Fatalf(packErrorFormat, fmt.Sprintf("pack action %s not found", action))
		}
	}
}

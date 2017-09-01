package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const usage = `
gopack: generate packs and tasks

usage: gopack generate --type=[pack|task] --name=<string> file
`

var (
	generateCommand = flag.NewFlagSet("generate", flag.ExitOnError)
	typeToGenerate  = generateCommand.String("type", "task", "task or pack")
	typeName        = generateCommand.String("name", "", "name of generated task or pack")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		generateCommand.Parse(os.Args[2:])
		switch *typeToGenerate {
		case "task":
			if err := generateTask(*typeName, generateCommand.Arg(0)); err != nil {
				fmt.Fprint(os.Stderr, usage)
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		case "pack":
			fmt.Fprint(os.Stderr, "not implemented")
		default:
			fmt.Fprint(os.Stderr, usage)
			os.Exit(1)
		}
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

const taskTmplt = `{{- $receiver := .ReceiverName }} {{- $arg := .ReceiverArg -}}
package {{ .PackageName }}

import (
	"fmt"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

// {{$receiver}} ...
type {{$receiver}} struct {
	Prop1 string
	Prop2 string

	gopack.BaseTask
}

// Run initializes default property values and delegates to BaseTask RunActions method
func ({{$arg}} {{$receiver}}) Run(runActions ...action.Enum) bool {
	{{$arg}}.setDefaults()
	return {{$arg}}.RunActions(&{{$arg}}, {{$arg}}.registerActions(), runActions)
}

func ({{$arg}} {{$receiver}}) registerActions() action.Methods {
	return action.Methods{
		action.Create: {{$arg}}.create,
	}
}

func ({{$arg}} {{$receiver}}) setDefaults() {
	if {{$arg}}.Prop1 == "" {
		{{$arg}}.Prop1 = "default value"
	}
}

// String returns a string which identifies the task with it's property values
func ({{$arg}} {{$receiver}}) String() string {
	return fmt.Sprintf("{{$receiver|ToLower}} %s %s", {{$arg}}.Prop1, {{$arg}}.Prop2)
}

func ({{$arg}} {{$receiver}}) create() (bool, error) {
	return true, nil
}
`

func generateTask(name, path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(p); err != nil && !os.IsNotExist(err) {
		return err
	} else if !os.IsNotExist(err) && !confirm(fmt.Sprintf("Overwrite %s (y/n)? ", p)) {
		return err
	}
	basePath := filepath.Base(filepath.Dir(p))
	if name == "" {
		name = strings.Split(filepath.Base(p), ".")[0]
		name = fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), strings.ToLower(name[1:]))
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	x := template.Must(template.New("GenTask").Funcs(funcMap).Parse(taskTmplt))
	b := &bytes.Buffer{}
	if err := x.Execute(b,
		struct {
			ReceiverName string
			ReceiverArg  string
			PackageName  string
		}{
			name,
			strings.ToLower(name[0:1]),
			basePath,
		}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(p, b.Bytes(), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Task %s generated", p)
	return nil
}

func confirm(prompt string) bool {
	response := ""
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return strings.TrimSpace(strings.ToLower(response)) == "y"
}

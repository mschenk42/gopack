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

usage: gopack generate --type=[runpack|pack|task] --name=<string> path

`

var (
	generateCommand = flag.NewFlagSet("generate", flag.ExitOnError)
	typeToGenerate  = generateCommand.String("type", "task", "task, pack or runpack")
	typeName        = generateCommand.String("name", "", "name of generated task, pack or runpack (defaults to path's base dir or file name)")
)

func main() {
	if len(os.Args) < 3 {
		exitOnError(nil)
	}

	switch os.Args[1] {
	case "generate":
		generateCommand.Parse(os.Args[2:])
		switch *typeToGenerate {
		case "task":
			if err := generateTask(*typeName, generateCommand.Arg(0), false); err != nil {
				exitOnError(err)
			}
		case "pack":
			if err := generatePack(*typeName, generateCommand.Arg(0), false); err != nil {
				exitOnError(err)
			}

		case "runpack":
			if err := generateRunPack(*typeName, generateCommand.Arg(0), false); err != nil {
				exitOnError(err)
			}
		default:
			exitOnError(nil)
		}
	default:
		exitOnError(nil)
	}
}

func exitOnError(err error) {
	fmt.Fprint(os.Stderr, usage)
	generateCommand.PrintDefaults()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	os.Exit(1)
}

func generateRunPack(name, path string, force bool) error {
	projectDir, err := filepath.Abs(path + "-runpack")
	if err != nil {
		return err
	}
	fi, err := os.Stat(projectDir)
	switch {
	case err != nil && !os.IsNotExist(err):
		return err
	case fi != nil && !fi.IsDir():
		return fmt.Errorf("path %s is a file and not a directory", projectDir)
	case !os.IsNotExist(err) && !force && !confirm(fmt.Sprintf("Overwrite %s (y/n)? ", projectDir)):
		return nil
	}
	if name == "" {
		name = filepath.Base(path)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, name), 0755); err != nil {
		return err
	}

	importPath, err := filepath.Rel(os.Getenv("GOPATH"), projectDir)
	if err != nil {
		return err
	}
	parts := strings.Split(importPath, string(filepath.Separator))
	parts = append(parts, filepath.Base(path))
	if len(parts) > 1 {
		importPath = filepath.Join(parts[1:]...)
	}
	if err := generatePack(name, filepath.Join(projectDir, name, name), true); err != nil {
		return err
	}

	x := template.Must(template.New("GenRunPackMain").Parse(runPackMain))
	b := &bytes.Buffer{}
	if err := x.Execute(b,
		struct {
			PackageName string
			ImportPath  string
		}{
			name,
			importPath,
		}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(projectDir, "main.go"), b.Bytes(), 0644); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"ToTitle": strings.Title,
	}
	x = template.Must(template.New("RunPackReadme").Funcs(funcMap).Parse(runPackReadme))
	b = &bytes.Buffer{}
	if err := x.Execute(b,
		struct {
			PackageName string
			ImportPath  string
		}{
			name,
			importPath,
		}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(projectDir, "README.md"), b.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "runpack %s generated\n", projectDir)
	return nil
}

func generatePack(name, path string, force bool) error {
	ext := filepath.Ext(path)
	if ext != ".go" {
		if ext != "" {
			path = strings.Replace(path, "."+ext, "", 1)
		}
		path += ".go"
	}
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(p)
	switch {
	case err != nil && !os.IsNotExist(err):
		return err
	case fi != nil && fi.IsDir():
		return fmt.Errorf("path %s is a directory and not a file", p)
	case !os.IsNotExist(err) && !force && !confirm(fmt.Sprintf("Overwrite %s (y/n)? ", p)):
		return nil
	}
	basePath := filepath.Base(filepath.Dir(p))
	if name == "" {
		name = strings.Split(filepath.Base(p), ".")[0]
		name = fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), strings.ToLower(name[1:]))
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	x := template.Must(template.New("GenPack").Funcs(funcMap).Parse(packTmplt))
	b := &bytes.Buffer{}
	if err := x.Execute(b,
		struct {
			PackageName string
		}{
			basePath,
		}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(p, b.Bytes(), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "pack %s generated\n", p)
	return nil
}

func generateTask(name, path string, force bool) error {
	ext := filepath.Ext(path)
	if ext != ".go" {
		if ext != "" {
			path = strings.Replace(path, "."+ext, "", 1)
		}
		path += ".go"
	}
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(p)
	switch {
	case err != nil && !os.IsNotExist(err):
		return err
	case fi != nil && fi.IsDir():
		return fmt.Errorf("path %s is a directory and not a file", p)
	case !os.IsNotExist(err) && !force && !confirm(fmt.Sprintf("Overwrite %s (y/n)? ", p)):
		return nil
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
	fmt.Fprintf(os.Stdout, "task %s generated\n", p)
	return nil
}

func confirm(prompt string) bool {
	response := ""
	fmt.Print(prompt)
	fmt.Scanln(&response)
	return strings.TrimSpace(strings.ToLower(response)) == "y"
}

const (
	runPackMain = `package main

import (
	"github.com/mschenk42/gopack"
	"{{.ImportPath}}"
)

func main() {
	{{ .PackageName }}.Run(gopack.LoadProperties())
}`
	runPackReadme = `# {{.PackageName|ToTitle}} Runpack
`
	packTmplt = `package {{ .PackageName }}

import (
	"fmt"

	"github.com/mschenk42/gopack"
)

// Run initializes the properties and runs the pack
func Run(props *gopack.Properties, actions []string) {
	pack := gopack.Pack{
		Name: "{{.PackageName}}",
		Props: &gopack.Properties{
			"{{.PackageName}}.prop1": "value",
		},
		Actions: actions,
		RunFunc: setup,
	}
	pack.Run(props)
}

var (
	actionMap = map[string]func(p *gopack.Pack){
		"default": run,
	}
)

func setup(pack *gopack.Pack) {
	if len(pack.Actions) == 0 {
		actionMap["default"](pack)
	}
	for _, action := range pack.Actions {
		if f, found := actionMap[action]; found {
			f(pack)
		} else {
			gopack.Log.Fatalf(gopack.PackErrorFormat, fmt.Sprintf("pack action %s not found", action))
		}
	}
}

func run(pack *gopack.Pack) {
	// run tasks and other packs within this method
}`
	taskTmplt = `{{- $receiver := .ReceiverName }} {{- $arg := .ReceiverArg -}}
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
func ({{$arg}} {{$receiver}}) Run(runActions ...action.Enum) gopack.ActionRunStatus {
	{{$arg}}.setDefaults()
	return {{$arg}}.RunActions(&{{$arg}}, {{$arg}}.registerActions(), runActions)
}

func ({{$arg}} {{$receiver}}) registerActions() action.Methods {
	return action.Methods{
		action.Create: {{$arg}}.create,
	}
}

func ({{$arg}} *{{$receiver}}) setDefaults() {
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
}`
)

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	packTmplt = `package main

import (
	"github.com/mschenk42/gopack"
)

func main() {
	props, actions := gopack.CLI()
	pack := gopack.Pack{
		Name: "{{.Name}}",
		Props: &gopack.Properties{
			"{{.Name}}.prop1": "val1",
		},
		Redact:  []string{"{{.Name}}.password"},
		Actions: actions,
		ActionMap: map[string]func(p *gopack.Pack){
			"default": run,
		},
	}
	pack.Run(props)
}

func run(pack *gopack.Pack) {

}`
	packReadmeTmplt = `# {{.Name|ToTitle}} Pack
`
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
func ({{$arg}} {{$receiver}}) Run(runActions ...action.Name) gopack.ActionRunStatus {
	{{$arg}}.setDefaults()
	return {{$arg}}.RunActions(&{{$arg}}, {{$arg}}.registerActions(), runActions)
}

func ({{$arg}} {{$receiver}}) registerActions() action.Funcs {
	return action.Funcs {
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

func generatePack(name, path string, force bool) error {
	projectDir, err := filepath.Abs(path + "-pack")
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
		name = strings.ToLower(filepath.Base(path))
	}
	if err := os.MkdirAll(filepath.Join(projectDir, name), 0755); err != nil {
		return err
	}

	x := template.Must(template.New("GenPack").Parse(packTmplt))
	b := &bytes.Buffer{}
	if err := x.Execute(b, struct{ Name string }{name}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(projectDir, "main.go"), b.Bytes(), 0644); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"ToTitle": strings.Title,
	}
	x = template.Must(template.New("PackReadme").Funcs(funcMap).Parse(packReadmeTmplt))
	b = &bytes.Buffer{}
	if err := x.Execute(b, struct{ Name string }{name}); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(projectDir, "README.md"), b.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "pack %s generated\n", projectDir)
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
	}
	name = fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), strings.ToLower(name[1:]))

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

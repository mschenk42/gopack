package task

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
)

type Template struct {
	Name   string
	Source string
	Props  *gopack.Properties
	Path   string
	Owner  string
	Group  string
	Mode   os.FileMode

	gopack.BaseTask
}

func (t Template) Run(runActions ...action.Enum) bool {
	t.setDefaults()
	return t.RunActions(&t, t.registerActions(), runActions)
}

func (t Template) registerActions() action.Methods {
	return action.Methods{
		action.Create: t.create,
	}
}

func (t *Template) setDefaults() {
	if t.Mode == 0 {
		t.Mode = 0755
	}
}

func (t Template) String() string {
	return fmt.Sprintf("template %s %s %s %s %s", t.Name, t.Path, t.Owner, t.Group, t.Mode)
}

func (t Template) create() (bool, error) {
	var (
		err          error
		chgTemplate  bool
		chgOwnership bool
		chgMode      bool
		fileExists   bool
		checkSumDiff bool
		fi           os.FileInfo
	)

	x := template.New(t.Name)
	if x, err = x.Parse(t.Source); err != nil {
		return false, gopack.NewTaskError(t, action.Create, err)
	}
	bt := &bytes.Buffer{}
	if err = x.Execute(bt, t.Props); err != nil {
		return false, gopack.NewTaskError(t, action.Create, err)
	}
	if fi, fileExists, err = fexists(t.Path); err != nil {
		return false, gopack.NewTaskError(t, action.Create, err)
	}
	if fileExists {
		bf := []byte{}
		if bf, err = ioutil.ReadFile(t.Path); err != nil {
			return false, gopack.NewTaskError(t, action.Create, err)
		}
		sumt := sha256.Sum256(bt.Bytes())
		sumf := sha256.Sum256(bf)
		checkSumDiff = sumt != sumf
	}
	if !fileExists || checkSumDiff {
		if err = ioutil.WriteFile(t.Path, bt.Bytes(), t.Mode); err != nil {
			return false, gopack.NewTaskError(t, action.Create, err)
		}
		chgTemplate = true
	} else {
		if fi.Mode().Perm() != t.Mode.Perm() {
			os.Chmod(t.Path, t.Mode)
			chgMode = true
		}
	}

	// do we need to set ownership?
	if t.Owner == "" && t.Group == "" {
		return chgTemplate || chgOwnership || chgMode, nil
	}
	if chgOwnership, err = chown(t.Path, t.Owner, t.Group); err != nil {
		return chgTemplate || chgOwnership || chgMode, gopack.NewTaskError(t, action.Create, err)
	} else {
		return chgTemplate || chgOwnership || chgMode, nil
	}
}

package task

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mschenk42/gopack"
)

type Template struct {
	Name   string
	Source string
	Path   string
	Owner  string
	Group  string
	Mode   os.FileMode

	props *gopack.Properties
	gopack.BaseTask
}

func (t Template) Run(props *gopack.Properties, runActions ...gopack.Action) bool {
	t.props = props
	t.setDefaults()
	return t.BaseTask.RunActions(&t, t.registerActions(), runActions)
}

func (t Template) Properties() *gopack.Properties {
	return t.props
}

func (t Template) registerActions() gopack.ActionMethods {
	return gopack.ActionMethods{
		gopack.CreateAction: t.create,
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
		fileExists   bool
		checkSumDiff bool
	)

	x := template.New(t.Name)
	if x, err = x.Parse(t.Source); err != nil {
		return false, t.Errorf(t, gopack.CreateAction, err)
	}
	bt := &bytes.Buffer{}
	if err = x.Execute(bt, t.props); err != nil {
		return false, t.Errorf(t, gopack.CreateAction, err)
	}
	if fileExists, err = fexists(t.Path); err != nil {
		return false, t.Errorf(t, gopack.CreateAction, err)
	}
	if fileExists {
		bf := []byte{}
		if bf, err = ioutil.ReadFile(t.Path); err != nil {
			return false, t.Errorf(t, gopack.CreateAction, err)
		}
		sumt := sha256.Sum256(bt.Bytes())
		sumf := sha256.Sum256(bf)
		checkSumDiff = sumt != sumf
	}
	if !fileExists || checkSumDiff {
		if err = ioutil.WriteFile(t.Path, bt.Bytes(), t.Mode); err != nil {
			return false, t.Errorf(t, gopack.CreateAction, err)
		}
		chgTemplate = true
	}

	// do we need to set ownership?
	if t.Owner == "" && t.Group == "" {
		return chgTemplate || chgOwnership, nil
	}
	if chgOwnership, err = chown(t.Path, t.Owner, t.Group); err != nil {
		return chgTemplate || chgOwnership, t.Errorf(t, gopack.CreateAction, err)
	} else {
		return chgTemplate || chgOwnership, nil
	}
}

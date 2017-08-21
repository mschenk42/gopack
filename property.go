package gopack

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Properties map[string]interface{}

func (p *Properties) Merge(props *Properties) {
	if props == nil {
		return
	}
	for k, v := range *props {
		(*p)[k] = v
	}
}

func (p *Properties) Load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return p.unmarshalJSON(b)
}

func (p *Properties) Save(filename string, x os.FileMode) error {
	b, err := p.marshalJSON()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, x)
}

func (p *Properties) unmarshalJSON(b []byte) error {
	return json.Unmarshal(b, p)
}

func (p *Properties) marshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

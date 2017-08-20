package gopack

import (
	"encoding/json"
	"fmt"
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

func (p *Properties) Bool() bool {
	return false
}

func (p *Properties) Float() float64 {
	return 0
}

func (p *Properties) Str(key string) (string, bool) {
	v, f := (*p)[key]
	switch v.(type) {
	case string:
		return v.(string), f
	default:
		panic(fmt.Sprintf("unable to convert %s to string", key))
	}
}

func (p *Properties) StrSlice() []string {
	return []string{}
}

func (p *Properties) StrMap() map[string]string {
	return map[string]string{}
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

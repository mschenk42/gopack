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

func (p *Properties) StringVar(key string) (string, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(string), found
	}
	return "", found
}

func (p *Properties) FloatVar(key string) (float64, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(float64), found
	}
	return 0, found
}

func (p *Properties) IntVar(key string) (int, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return int(v.(float64)), found
	}
	return 0, found
}

func (p *Properties) MapVar(key string) (map[string]interface{}, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(map[string]interface{}), found
	}
	return map[string]interface{}{}, found
}

func (p *Properties) SliceVar(key string) ([]interface{}, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.([]interface{}), found
	}
	return []interface{}{}, found
}

func (p *Properties) BoolVar(key string) (bool, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(bool), found
	}
	return false, found
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

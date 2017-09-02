package gopack

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Properties map[string]interface{}

func (p *Properties) String() string {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		Log.Fatalf("    ! %s", err)
	}
	return string(b)
}

func (p Properties) Redact(redact []string) *Properties {
	for _, k := range redact {
		p[k] = "***"
	}
	return &p
}

func (p *Properties) Merge(props *Properties) {
	if props == nil {
		return
	}
	for k, v := range *props {
		(*p)[k] = v
	}
}

func (p *Properties) Str(key string) (string, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(string), found
	}
	return "", found
}

func (p *Properties) Float(key string) (float64, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(float64), found
	}
	return 0, found
}

func (p *Properties) Int(key string) (int, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return int(v.(float64)), found
	}
	return 0, found
}

func (p *Properties) Map(key string) (map[string]interface{}, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(map[string]interface{}), found
	}
	return map[string]interface{}{}, found
}

func (p *Properties) Slice(key string) ([]interface{}, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.([]interface{}), found
	}
	return []interface{}{}, found
}

func (p *Properties) Bool(key string) (bool, bool) {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(bool), found
	}
	return false, found
}

func (p *Properties) Load(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.unmarshalJSON(b)
}

func (p *Properties) Save(w io.Writer) error {
	b, err := p.marshalJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (p *Properties) unmarshalJSON(b []byte) error {
	return json.Unmarshal(b, p)
}

func (p *Properties) marshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

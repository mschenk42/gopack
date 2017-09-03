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

func (p *Properties) Exists(key string) bool {
	_, x := (*p)[key]
	return x
}

func (p *Properties) Str(key string) string {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(string)
	}
	return ""
}

func (p *Properties) Float(key string) float64 {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(float64)
	}
	return 0
}

func (p *Properties) Int(key string) int {
	v, found := (*p)[key]
	if found && v != nil {
		return int(v.(float64))
	}
	return 0
}

func (p *Properties) Map(key string) map[string]interface{} {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (p *Properties) Slice(key string) []interface{} {
	v, found := (*p)[key]
	if found && v != nil {
		return v.([]interface{})
	}
	return []interface{}{}
}

func (p *Properties) Bool(key string) bool {
	v, found := (*p)[key]
	if found && v != nil {
		return v.(bool)
	}
	return false
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

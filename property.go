package gopack

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/mschenk42/gopack/color"
)

type Properties map[string]interface{}

func (p *Properties) String() string {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		Log.Fatalf(color.Red("! %s"), err)
	}
	return string(b)
}

func (p Properties) Redact(redact []string) *Properties {
	pcopy := Properties{}
	for k, v := range p {
		pcopy[k] = v
	}
	for _, k := range redact {
		pcopy[k] = "***"
	}
	return &pcopy
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
	if !found || v == nil {
		return ""
	}
	x, _ := v.(string)
	return x
}

func (p *Properties) StrRequired(key string) string {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to string", key))
	}
	x, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to string", key))
	}
	return x
}

func (p *Properties) Float(key string) float64 {
	v, found := (*p)[key]
	if !found || v == nil {
		return 0
	}
	x, _ := v.(float64)
	return x
}

func (p *Properties) FloatRequired(key string) float64 {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to float", key))
	}
	x, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to float", key))
	}
	return x
}

func (p *Properties) Int(key string) int {
	v, found := (*p)[key]
	if !found || v == nil {
		return 0
	}
	x, _ := v.(float64)
	return int(x)
}

func (p *Properties) IntRequired(key string) int {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to int", key))
	}
	x, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to int", key))
	}
	return int(x)
}

func (p *Properties) Map(key string) map[string]interface{} {
	v, found := (*p)[key]
	if !found || v == nil {
		return map[string]interface{}{}
	}
	x, _ := v.(map[string]interface{})
	return x
}

func (p *Properties) MapRequired(key string) map[string]interface{} {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to map", key))
	}
	x, ok := v.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to map", key))
	}
	return x
}

func (p *Properties) Slice(key string) []interface{} {
	v, found := (*p)[key]
	if !found || v == nil {
		return []interface{}{}
	}
	x, _ := v.([]interface{})
	return x
}

func (p *Properties) SliceRequired(key string) []interface{} {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to slice", key))
	}
	x, ok := v.([]interface{})
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to slice", key))
	}
	return x
}

func (p *Properties) Bool(key string) bool {
	v, found := (*p)[key]
	if !found || v == nil {
		return false
	}
	x, _ := v.(bool)
	return x
}

func (p *Properties) BoolRequired(key string) bool {
	v, found := (*p)[key]
	if !found || v == nil {
		panic(fmt.Sprintf("unable to convert %s to boolean", key))
	}
	x, ok := v.(bool)
	if !ok {
		panic(fmt.Sprintf("unable to convert %s to boolean", key))
	}
	return x
}

func (p *Properties) Load(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.unmarshalJSON(b)
}

func (p *Properties) Save(w io.Writer) error {
	// let's save it in a pretty json format
	_, err := w.Write([]byte(p.String()))
	return err
}

func (p *Properties) unmarshalJSON(b []byte) error {
	return json.Unmarshal(b, p)
}

func (p *Properties) marshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

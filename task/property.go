package task

import "fmt"

type Properties map[string]interface{}

func (p *Properties) Merge(x Properties) {
	for k, v := range x {
		(*p)[k] = v
	}
}

func (p *Properties) Bool() bool {
	return false
}

func (p *Properties) Float() float64 {
	return 0
}

func (p *Properties) Str(k string) string {
	v := (*p)[k]
	switch v.(type) {
	case string:
		return v.(string)
	default:
		handleError(fmt.Errorf("unable to convert %s to string", k))
	}
	return ""
}

func (p *Properties) StrSlice() []string {
	return []string{}
}

func (p *Properties) StrMap() map[string]string {
	return map[string]string{}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

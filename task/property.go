package task

import "fmt"

type Properties map[string]interface{}

func (p *Properties) Merge(props Properties) {
	for k, v := range props {
		(*p)[k] = v
	}
}

func (p *Properties) Bool() bool {
	return false
}

func (p *Properties) Float() float64 {
	return 0
}

func (p *Properties) Str(key string) string {
	v := (*p)[key]
	switch v.(type) {
	case string:
		return v.(string)
	default:
		handleError(fmt.Errorf("unable to convert %s to string", key))
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

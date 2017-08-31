package gopack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"nginx.log_dir":"/var/log/nginx", "nginx.cache_dir": "/var/cache"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal("/var/log/nginx", p["nginx.log_dir"])
	assert.Equal("/var/cache", p["nginx.cache_dir"])
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	b := []byte(`{"nginx.log_dir":"/var/log/nginx", "nginx.cache_dir": "/var/cache"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)

	b = []byte(`{"nginx.cache_dir": "/etc/cache"}`)
	override := &Properties{}
	err = p.unmarshalJSON(b)
	assert.NoError(err)

	assert.NotPanics(func() { p.Merge(override) })
	assert.Equal(Properties{"nginx.log_dir": "/var/log/nginx", "nginx.cache_dir": "/etc/cache"}, p)
}

func TestStringType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"string": "val"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Str("string")
	assert.Equal(true, f)
	assert.Equal("val", v)
}

func TestFloat64Type(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"float":0.1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Float("float")
	assert.Equal(true, f)
	assert.Equal(0.1, v)
}

func TestIntType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"int": 1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Int("int")
	assert.Equal(true, f)
	assert.Equal(1, v)
}

func TestMapType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"map": {"key1": "val1", "key2": "val2"}}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Map("map")
	assert.Equal(true, f)
	assert.Equal(map[string]interface{}{"key1": "val1", "key2": "val2"}, v)
}

func TestSliceType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"array": ["val1", "val2"]}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Slice("array")
	assert.Equal(true, f)
	assert.Equal([]interface{}{"val1", "val2"}, v)
}

func TestBoolType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"bool": true}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Bool("bool")
	assert.Equal(true, f)
	assert.Equal(true, v)
}

func TestNilType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"nil": null }`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	v, f := p.Map("nil")
	assert.Equal(true, f)
	assert.Equal(map[string]interface{}{}, v)
}

func TestWrongType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"string": "val"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Panics(func() { fmt.Println("invalid type assertion, shouldn't print", p["string"].(float64)) })
}

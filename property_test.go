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
	assert.Equal(true, p.Exists("string"))
	assert.Equal("val", p.Str("string"))
	assert.NotPanics(func() { p.StrRequired("string") })
	assert.Panics(func() { p.StrRequired("invalid") })
}

func TestFloat64Type(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"float":0.1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("float"))
	assert.Equal(0.1, p.Float("float"))
	assert.NotPanics(func() { p.FloatRequired("float") })
	assert.Panics(func() { p.FloatRequired("invalid") })
}

func TestIntType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"int": 1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("int"))
	assert.Equal(1, p.Int("int"))
	assert.NotPanics(func() { p.IntRequired("int") })
	assert.Panics(func() { p.IntRequired("invalid") })
}

func TestMapType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"map": {"key1": "val1", "key2": "val2"}}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("map"))
	assert.Equal(map[string]interface{}{"key1": "val1", "key2": "val2"}, p.Map("map"))
	assert.NotPanics(func() { p.MapRequired("map") })
	assert.Panics(func() { p.MapRequired("invalid") })
}

func TestSliceType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"array": ["val1", "val2"]}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("array"))
	assert.Equal([]interface{}{"val1", "val2"}, p.Slice("array"))
	assert.NotPanics(func() { p.SliceRequired("array") })
	assert.Panics(func() { p.SliceRequired("invalid") })
}

func TestBoolType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"bool": true}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("bool"))
	assert.Equal(true, p.Bool("bool"))
	assert.NotPanics(func() { p.BoolRequired("bool") })
	assert.Panics(func() { p.BoolRequired("invalid") })
}

func TestNilType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"nil": null }`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p.Exists("nil"))
	assert.Equal(map[string]interface{}{}, p.Map("nil"))
}

func TestWrongType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"string": "val"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Panics(func() { fmt.Println("invalid type assertion, shouldn't print", p["string"].(float64)) })
}

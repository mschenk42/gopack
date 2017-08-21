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
	assert.Equal("val", p["string"].(string))
}

func TestFloat64Type(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"float":0.1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(0.1, p["float"].(float64))
}

func TestIntType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"int": 1}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(1, int(p["int"].(float64)))
}

func TestMapType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"map": {"key1": "val1", "key2": "val2"}}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(map[string]interface{}{"key1": "val1", "key2": "val2"}, p["map"].(map[string]interface{}))
}

func TestArrayType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"array": ["val1", "val2"]}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal([]interface{}{"val1", "val2"}, p["array"].([]interface{}))
}

func TestBoolType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"bool": true}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(true, p["bool"].(bool))
}

func TestNilType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"nil": null }`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Equal(nil, p["nil"])
}

func TestWrongType(t *testing.T) {
	assert := assert.New(t)
	b := []byte(`{"string": "val"}`)
	p := Properties{}
	err := p.unmarshalJSON(b)
	assert.NoError(err)
	assert.Panics(func() { fmt.Println("invalid type assertion, shouldn't print", p["string"].(float64)) })
}

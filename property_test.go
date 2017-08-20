package gopack

import (
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

package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/stretchr/testify/assert"
)

func TestCreateTemplate(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/test-create-template"

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer func() { os.RemoveAll(testDir) }()

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir: {{ index . "nginx.log_dir"}}`,
	}.Run(
		&gopack.Properties{"nginx.log_dir": "/var/log/nginx"},
		gopack.CreateAction,
	)

	b, err := ioutil.ReadFile(fmt.Sprintf("%s/mypack.conf", testDir))
	assert.NoError(err)
	assert.Regexp(`log_dir: /var/log/nginx`, string(b))
}

func TestUpToDateTemplate(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/test-uptodate-template"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer func() { os.RemoveAll(testDir) }()

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir: {{ index . "nginx.log_dir"}}`,
	}.Run(
		&gopack.Properties{"nginx.log_dir": "/var/log/nginx"},
		gopack.CreateAction,
	)

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir: {{ index . "nginx.log_dir"}}`,
	}.Run(
		&gopack.Properties{"nginx.log_dir": "/var/log/nginx"},
		gopack.CreateAction,
	)

	assert.Regexp(`.*template mypack /tmp/.*/mypack.conf.*create.*(up to date).*`, buf.String())
	fmt.Print(buf.String())
}

func TestNotUpToDateTemplate(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/test-not-uptodate-template"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer func() { os.RemoveAll(testDir) }()

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir: {{ index . "nginx.log_dir"}}`,
	}.Run(
		&gopack.Properties{"nginx.log_dir": "/var/log/nginx"},
		gopack.CreateAction,
	)

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir: {{ index . "nginx.log_dir"}}`,
	}.Run(
		&gopack.Properties{"nginx.log_dir": "/log/nginx"},
		gopack.CreateAction,
	)

	assert.NotRegexp(`.*template mypack /tmp/.*/mypack.conf.*create.*(up to date).*`, buf.String())
	fmt.Print(buf.String())
}

func TestModeUpdateTemplate(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/test-mode-update-template"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer func() { os.RemoveAll(testDir) }()

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0755,
		Source: `log_dir:`,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	assert.Regexp(`.*-rwxr-xr-x:.*`, buf.String())

	Template{
		Name:   "mypack",
		Path:   fmt.Sprintf("%s/mypack.conf", testDir),
		Mode:   0775,
		Source: `log_dir:`,
	}.Run(
		nil,
		gopack.CreateAction,
	)

	assert.Regexp(`.*-rwxrwxr-x:.*`, buf.String())
	assert.NotRegexp(`.*template mypack /tmp/.*/mypack.conf.*create.*(up to date).*`, buf.String())
	fmt.Print(buf.String())
}

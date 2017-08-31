package task

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/mschenk42/gopack/action"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	assert := assert.New(t)
	const testDir = "/tmp/create-dir"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(action.Create)
	defer os.Remove(testDir)

	_, err := os.Stat(testDir)
	assert.Nil(err)
	assert.Regexp(`.*directory.*/tmp/create-dir.*create.*(run)`, buf.String())
}

func TestCreateExistingDirectory(t *testing.T) {
	assert := assert.New(t)
	const testDir = "/tmp/create-existing-dir"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	err := os.Mkdir(testDir, 0755)
	defer os.Remove(testDir)
	assert.Nil(err)

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(action.Create)

	_, err = os.Stat(testDir)
	assert.Nil(err)
	assert.Regexp(`.*directory.*/tmp/create-existing-dir.*create.*(up to date)`, buf.String())
}

func TestCreateDirectoryValidOwner(t *testing.T) {
	// assert := assert.New(t)
	// const testDir = "/tmp/create-dir-owner"

	// Directory{
	// 	Owner: "mschenk",
	// 	Group: "admin",
	// 	Path:  testDir,
	// 	Mode:  0755,
	// }.Run(action.CreateAction)
	// defer os.Remove(testDir)

	// _, err := os.Stat(testDir)
	// assert.Nil(err)

}

func TestRemoveDirectory(t *testing.T) {
	assert := assert.New(t)
	const testDir = "/tmp/remove-dir"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	err := os.Mkdir(testDir, 0755)
	defer os.Remove(testDir)
	assert.Nil(err)

	Directory{
		Path: testDir,
	}.Run(action.Remove)

	_, err = os.Stat(testDir)
	assert.NotNil(err)
	assert.Regexp(`.*directory.*/tmp/remove-dir.*remove.*(run)`, buf.String())
}

func TestRemoveMissingDirectory(t *testing.T) {
	assert := assert.New(t)
	const testDir = "/tmp/remove-missing-dir"

	saveLogger := gopack.Log
	buf := &bytes.Buffer{}
	gopack.Log = log.New(buf, "", 0)
	defer func() { gopack.Log = saveLogger }()

	_, err := os.Stat(testDir)
	assert.NotNil(err)

	Directory{
		Path: testDir,
	}.Run(action.Remove)

	_, err = os.Stat(testDir)
	assert.NotNil(err)
	assert.Regexp(`.*directory.*/tmp/remove-missing-dir.*remove.*(up to date)`, buf.String())
}

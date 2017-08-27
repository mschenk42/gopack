package task

import (
	"os"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/create_dir"
	Directory{
		Path: d,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer os.Remove(d)

	_, err := os.Stat(d)
	assert.Nil(err)
}

func TestCreateExistingDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/create_existing_dir"
	err := os.Mkdir(d, 0755)
	defer os.Remove(d)
	assert.Nil(err)

	Directory{
		Path: d,
		Mode: 0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)

	_, err = os.Stat(d)
	assert.Nil(err)
}

func TestCreateDirectoryValidOwner(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/create_dir_owner"

	Directory{
		Owner: "mschenk",
		Group: "admin",
		Path:  d,
		Mode:  0755,
	}.Run(
		nil,
		gopack.CreateAction,
	)
	defer os.Remove(d)

	_, err := os.Stat(d)
	assert.Nil(err)

}

func TestRemoveDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/remove_dir"
	err := os.Mkdir(d, 0755)
	defer os.Remove(d)
	assert.Nil(err)

	Directory{
		Path: d,
	}.Run(
		nil,
		gopack.RemoveAction,
	)

	_, err = os.Stat(d)
	assert.NotNil(err)
}

func TestRemoveMissingDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/remove_missing_dir"
	_, err := os.Stat(d)
	assert.NotNil(err)

	Directory{
		Path: d,
	}.Run(
		nil,
		gopack.RemoveAction,
	)

	_, err = os.Stat(d)
	assert.NotNil(err)
}

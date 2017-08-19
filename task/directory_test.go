package task

import (
	"log"
	"os"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/test"
	Directory{
		Path: d,
		Perm: 0755,
	}.Run(
		nil,
		log.New(os.Stdout, "", 0),
		gopack.CreateAction,
	)
	defer os.Remove(d)

	_, err := os.Stat(d)
	assert.Nil(err)
}

func TestCreateExistingDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/test"
	err := os.Mkdir(d, 0755)
	defer os.Remove(d)
	assert.Nil(err)

	Directory{
		Path: d,
		Perm: 0755,
	}.Run(
		nil,
		log.New(os.Stdout, "", 0),
		gopack.CreateAction,
	)

	_, err = os.Stat(d)
	assert.Nil(err)
}

func TestRemoveDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/test"
	err := os.Mkdir(d, 0755)
	defer os.Remove(d)
	assert.Nil(err)

	Directory{
		Path: d,
	}.Run(
		nil,
		log.New(os.Stdout, "", 0),
		gopack.RemoveAction,
	)

	_, err = os.Stat(d)
	assert.NotNil(err)
}

func TestRemoveMissingDirectory(t *testing.T) {
	assert := assert.New(t)

	d := "/tmp/test"
	_, err := os.Stat(d)
	assert.NotNil(err)

	Directory{
		Path: d,
	}.Run(
		nil,
		log.New(os.Stdout, "", 0),
		gopack.RemoveAction,
	)

	_, err = os.Stat(d)
	assert.NotNil(err)
}

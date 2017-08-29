package task

import (
	"os"
	"testing"

	"github.com/mschenk42/gopack"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/create-dir"

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(gopack.CreateAction)
	defer os.Remove(testDir)

	_, err := os.Stat(testDir)
	assert.Nil(err)
}

func TestCreateExistingDirectory(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/create-existing-dir"

	err := os.Mkdir(testDir, 0755)
	defer os.Remove(testDir)
	assert.Nil(err)

	Directory{
		Path: testDir,
		Mode: 0755,
	}.Run(gopack.CreateAction)

	_, err = os.Stat(testDir)
	assert.Nil(err)
}

func TestCreateDirectoryValidOwner(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/create-dir-owner"

	Directory{
		Owner: "mschenk",
		Group: "admin",
		Path:  testDir,
		Mode:  0755,
	}.Run(gopack.CreateAction)
	defer os.Remove(testDir)

	_, err := os.Stat(testDir)
	assert.Nil(err)

}

func TestRemoveDirectory(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/remove-dir"

	err := os.Mkdir(testDir, 0755)
	defer os.Remove(testDir)
	assert.Nil(err)

	Directory{
		Path: testDir,
	}.Run(gopack.RemoveAction)

	_, err = os.Stat(testDir)
	assert.NotNil(err)
}

func TestRemoveMissingDirectory(t *testing.T) {
	assert := assert.New(t)

	const testDir = "/tmp/remove-missing-dir"

	_, err := os.Stat(testDir)
	assert.NotNil(err)

	Directory{
		Path: testDir,
	}.Run(gopack.RemoveAction)

	_, err = os.Stat(testDir)
	assert.NotNil(err)
}

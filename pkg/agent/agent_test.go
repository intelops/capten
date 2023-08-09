package agent

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAppConfigs(t *testing.T) {
	assert := require.New(t)

	dir, err := os.Getwd()
	assert.Nil(err)

	finalPath := path.Join(dir, "../..", "/apps/conf")
	confs, err := readAppConfigs(finalPath)
	assert.Nil(err)
	for _, c := range confs {
		_, err := c.ToSyncAppData()
		assert.Nil(err)
	}

	files, err := ioutil.ReadDir(finalPath)
	assert.Equal(len(files), len(confs))
	assert.Nil(err)
}

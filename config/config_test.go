package config_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"

	"github.com/rockset/cli/config"
	"github.com/stretchr/testify/require"
)

func tmpDir(t *testing.T) (string, func()) {
	tmp, err := os.MkdirTemp(os.TempDir(), "home")
	require.NoError(t, err)

	return tmp, func() {
		assert.NoError(t, os.RemoveAll(tmp))
	}
}

func TestConfigFile(t *testing.T) {
	assert.NotEmpty(t, config.File)
	t.Logf("config.File: %s", config.File)
}

func TestLoadFile_noFile(t *testing.T) {
	_, err := config.LoadFile(path.Join("testdata", "missing.yaml"))
	require.NoError(t, err)
}

func TestStore_noFile(t *testing.T) {
	dir, cleanup := tmpDir(t)
	defer cleanup()

	cfg := config.New()
	file := path.Join(dir, "foo.yaml")
	err := config.StoreFile(cfg, file)
	require.NoError(t, err)

	f, err := os.ReadFile(file)
	assert.Equal(t, `current: ""
keys: {}
tokens: {}
`, string(f))
}

func TestConfig(t *testing.T) {
	cfg := config.New()
	cfg.Keys["foo"] = config.APIKey{}
	cfg.Tokens["bar"] = config.Token{}
}

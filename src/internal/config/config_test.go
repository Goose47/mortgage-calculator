package config

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func setup(t *testing.T, cfg *Config) (*os.File, func()) {
	data, err := yaml.Marshal(&cfg)
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp("", "config-*.yml")
	require.NoError(t, err)

	_, err = tmpFile.Write(data)
	require.NoError(t, err)

	cleanup := func() {
		tmpFile.Close()
	}

	return tmpFile, cleanup
}

func TestFetchConfigPath(t *testing.T) {
	require.Equal(t, fetchConfigPath(), os.Getenv("CONFIG_PATH")) // no --config flag in tests
}

func TestLoadPath(t *testing.T) {
	cfg := &Config{
		Env:  "local",
		Port: 8080,
		Cache: Cache{
			TTL:   100,
			Clear: 100,
		},
	}

	file, cleanup := setup(t, cfg)
	defer cleanup()

	res, err := LoadPath(file.Name())
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, *cfg, *res)

	res, err = LoadPath("")
	require.Error(t, err)
	require.Empty(t, res)
}

func TestMustLoadPath(t *testing.T) {
	cfg := &Config{
		Env:  "local",
		Port: 8080,
		Cache: Cache{
			TTL:   100,
			Clear: 100,
		},
	}

	file, cleanup := setup(t, cfg)
	defer cleanup()

	res := MustLoadPath(file.Name())
	require.NotEmpty(t, res)
	require.Equal(t, *cfg, *res)

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}()
	MustLoadPath("")
}

func TestMustLoad(t *testing.T) {
	cfg := &Config{
		Env:  "local",
		Port: 8080,
		Cache: Cache{
			TTL:   100,
			Clear: 100,
		},
	}

	file, cleanup := setup(t, cfg)
	defer cleanup()

	err := os.Setenv("CONFIG_PATH", file.Name())
	require.NoError(t, err)

	res := MustLoad()
	require.NotEmpty(t, res)
	require.Equal(t, *cfg, *res)

	err = file.Close()
	require.NoError(t, err)
	err = os.Remove(file.Name())
	require.NoError(t, err)

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}()
	MustLoadPath(file.Name())
}

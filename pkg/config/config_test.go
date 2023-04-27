package config_test

import (
	"os"
	"restapi/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Config struct {
	Mode string
	App  struct {
		Name string
		Port int64
	}
	Cache struct {
		Enable        bool
		ExpiredMinute int64 `mapstructure:"expired_minute"`
	}
}

var (
	configFile = "./testdata/full.yaml"
	// yamlFullConfig return the config associate with the configFile
	yamlFullConfig = func() Config {
		var c Config
		c.Mode = "sandbox"
		c.App.Name = "refund-core"
		c.App.Port = 8080
		c.Cache.Enable = true
		c.Cache.ExpiredMinute = 10
		return c
	}
)

func TestLoadNotExistedConfigFile(t *testing.T) {
	var c Config
	err := config.Load("./abc.yaml", &c)
	require.Error(t, err)
}

func TestLoadFullFromYaml(t *testing.T) {
	var c Config
	err := config.Load("./testdata/full.yaml", &c)
	require.NoError(t, err)

	assert.Equal(t, yamlFullConfig(), c)
}

func TestLoadFromEnv(t *testing.T) {
	setEnv(t, "APP_NAME", "refund")

	var c Config
	err := config.Load(configFile, &c)
	require.NoError(t, err)

	want := yamlFullConfig()
	want.App.Name = "refund"
	assert.Equal(t, want, c)
}

func TestLoadConfigHasUnderscoreFromEnv(t *testing.T) {
	setEnv(t, "CACHE_EXPIRED_MINUTE", "15")

	var c Config
	err := config.Load(configFile, &c)
	require.NoError(t, err)

	want := yamlFullConfig()
	want.Cache.ExpiredMinute = 15
	assert.Equal(t, want, c)
}

func setEnv(t *testing.T, k, v string) {
	originV := os.Getenv(k)
	t.Cleanup(func() {
		_ = os.Setenv(k, originV)
	})

	_ = os.Setenv(k, v)
}

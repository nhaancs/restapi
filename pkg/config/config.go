package config

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Load config into config, config must be a pointer to the config struct
//
// If the file is found, fully replace config with content in the file.
//
// Environment variable will override the config
// For example:
//
//	type Config struct {
//			Mode string
//			Cache struct {
//				Addr string
//				ExpiredMinute int `mapstructure:"expired_minute"`
//			}
//	}
//
// Field Config.Mode will be replaced with evn MODE
// Field Config.Cache.Addr  will be replaced with evn CACHE_ADDR
// Field Config.Cache.ExpiredMinute  will be replaced with evn CACHE_EXPIRED_MINUTE
func Load(file string, config interface{}) error {
	v := viper.New()
	m := make(map[string]interface{})

	if err := mapstructure.Decode(config, &m); err != nil {
		return fmt.Errorf("mapstructure: %v", err)
	}

	if err := v.MergeConfigMap(m); err != nil {
		return fmt.Errorf("merge config map: %v", err)
	}

	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("unmarshal config: %v", err)
	}

	return nil
}

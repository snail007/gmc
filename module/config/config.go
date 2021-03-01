// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gconfig

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	*viper.Viper
}

func (c *Config) Sub(key string) gcore.SubConfig {
	v := c.Viper.Sub(key)
	if reflect.ValueOf(v).IsNil() {
		return nil
	}
	return v
}

func NewConfig() *Config {
	c := &Config{
		Viper: viper.New(),
	}
	bindEnv(c)
	return c
}

func NewConfigFile(file string) (c *Config, err error) {
	cfg := viper.New()
	cfg.SetConfigFile(file)
	err = cfg.ReadInConfig()
	if err != nil {
		return
	}
	c = &Config{
		Viper: cfg,
	}
	bindEnv(c)
	return
}

func NewConfigBytes(b []byte) (c *Config, err error) {
	c = &Config{
		Viper: viper.New(),
	}
	c.SetConfigType("toml")
	bindEnv(c)
	err = c.ReadConfig(bytes.NewReader(b))
	return
}
func bindEnv(config gcore.Config) {
	// env binding
	prefix := os.Getenv("ENV_PREFIX")
	if prefix == "" {
		prefix = "GMC"
	}
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetEnvPrefix(prefix)
	config.AutomaticEnv()
}

// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gconfig

import (
	"bytes"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
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

func New() *Config {
	c := &Config{
		Viper: viper.New(),
	}
	bindEnv(c)
	return c
}
func NewFromSearch(paths []string, filename string, typ ...string) (c *Config, err error) {
	filePath := ""
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(absPath, filename)
		_, err = os.Stat(filePath)
		if err == nil {
			break
		} else {
			filePath = ""
		}
	}
	if filePath == "" {
		return nil, fmt.Errorf("file [%s] not found, in paths: %v", filename, paths) // 发生其他错误
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	c, err = NewConfigBytes(b, typ...)
	return
}

func NewFromFile(file string, typ ...string) (c *Config, err error) {
	cfg := viper.New()
	cfg.SetConfigFile(file)
	if len(typ) == 1 {
		cfg.SetConfigType(typ[0])
	}
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

func NewConfigBytes(b []byte, typ ...string) (c *Config, err error) {
	c = &Config{
		Viper: viper.New(),
	}
	if len(typ) == 1 {
		c.SetConfigType(typ[0])
	} else {
		c.SetConfigType("toml")
	}
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

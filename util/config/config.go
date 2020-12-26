package gconfig

import (
	"bytes"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	return &Config{
		Viper: viper.New(),
	}
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
	return
}

func NewConfigBytes(b []byte) *Config {
	c := &Config{
		Viper: viper.New(),
	}
	c.SetConfigType("toml")
	c.ReadConfig(bytes.NewReader(b))
	return c
}

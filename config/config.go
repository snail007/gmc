package gconfig

import (
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
func NewConfigFile(file string) (c *Config,err error){
	cfg:=viper.New()
	cfg.SetConfigFile(file)
	err=cfg.ReadInConfig()
	if err!=nil{
		return
	}
	c= &Config{
		Viper: cfg,
	}
	return
}
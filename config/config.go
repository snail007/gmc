package gmcconfig

import (
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}
func NewFile(file string) (c *Config ,err error){
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
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

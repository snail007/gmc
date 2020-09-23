package gmcconfig

import (
	"github.com/spf13/viper"
)

type GMCConfig struct {
	*viper.Viper
}

func New() *GMCConfig {
	return &GMCConfig{
		Viper: viper.New(),
	}
}

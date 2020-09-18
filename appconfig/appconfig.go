package appconfig

import (
	"github.com/snail007/gmc/session"
	"github.com/spf13/viper"
)

var (
	SessionStore session.Store
	Config       *viper.Viper
	isParsed     bool
)

func Parse() (err error) {
	if isParsed {
		return
	}
	defer func() {
		if err == nil {
			isParsed = true
		}
	}()
	Config = viper.New()
	Config.SetConfigName("app")
	Config.AddConfigPath(".")
	Config.AddConfigPath("config")
	Config.AddConfigPath("conf")
	Config.AddConfigPath("../appconfig")
	err = Config.ReadInConfig()
	if err != nil {
		return
	}
	initSessionStore()
	return
}

func initSessionStore() {

}

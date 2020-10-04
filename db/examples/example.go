package main

import (
	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.New.Config()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init only using [database] section in app.toml
	gmc.DB.Init(cfg)

	// database default is mysql in app.toml
	// so gmc.DB() equal to  gmc.DBMySQL()
	// we can connect to multiple cache drivers at same time, id is the unique name of driver
	// gmc.DB(id) to load `id` named default driver.
	db := gmc.DB.DB().(*gmc.MySQL)
	//do something with db
	db.AR()
}

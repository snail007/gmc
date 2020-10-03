package main

import (
	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.NewConfig()
	cfg.SetConfigFile("../../app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init only using [database] section in app.toml
	gmc.InitDB(cfg)

	// database default is mysql in app.toml
	// so gmc.DB() equal to  gmc.DBMySQL()
	// we can connect to multiple cache drivers at same time, id is the unique name of driver
	// gmc.DB(id) to load `id` named default driver.
	db := gmc.DB().(*gmc.MySQL)
	//do something with db
}

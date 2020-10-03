package main

import (
	"time"

	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.NewConfig()
	cfg.SetConfigFile("../../app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init only using [cache] section in app.toml
	gmc.InitCache(cfg)

	// cache default is redis in app.toml
	// so gmc.Cache() equal to  gmc.Redis()
	// we can connect to multiple cache drivers at same time, id is the unique name of driver
	// gmc.Cache(id) to load `id` named default driver.
	c := gmc.Cache()
	c.Set("test", "aaa", time.Second)
	c.Get("test")
}

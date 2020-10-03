# GMC Cache

1. Support of Redis.
1. Support of Multiple redis source.

## Configuration
cache configuration section in app.toml
```toml
########################################################
# cache configuration
########################################################
# redis supported
# support of mutiple redis server 
# notic: each config section must have an unique id 
########################################################
[cache]
default="redis"
[[cache.redis]]
debug=true
enable=true
id="default"
address="127.0.0.1:6379"
prefix=""
password=""
# seconds
timeout=10
dbnum=0
maxidle=10
maxactive=30
# seconds
idletimeout=300
# seconds
maxconnlifetime=3600
wait=true
```

## Example

```go
package main

import (
	"github.com/snail007/gmc"
	"time"
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
```
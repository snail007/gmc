# GMC DATABASE

1. Support of MYSQL , SQLITE3.
1. Support of Multiple database source.
1. Support of encrypt sqlite3 databse.

## Configuration
database configuration section in app.toml
```toml
########################################################
# database configuration
########################################################
# mysql,sqlite3 are both supported
# support of mutiple mysql server 
# support of mutiple sqlite3 database
# notic: each config section must have an unique id 
########################################################
[database]
default="mysql"
[[database.mysql]]
enable=true
id="default"
host="127.0.0.1"
port="3306"
username="root"
password="admin"
database="test"
prefix=""
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=3000
readtimeout=5000
writetimeout=5000
maxlifetimeseconds=1800

[[database.mysql]]
enable=false
id="news"
host="127.0.0.1"
port="3306"
username="root"
password="admin"
database="test"
prefix=""
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=3000
readtimeout=5000
writetimeout=5000
maxlifetimeseconds=1800

[[database.sqlite3]]
enable=false
id="default"
database="test.db"
# if password is not empty , database will be encrypted.
password=""
prefix=""
prefix_sql_holder="__PREFIX__"
# syncmode 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
syncmode=0
# openmode ro,rw,rwc,memory
openmode="rw"
# shared,private
cachemode="shared"
```

## Example

```go
package main

import (
	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.New.Config()
	cfg.SetConfigFile("../../app.toml")
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
}
```
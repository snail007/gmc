package gdb

import (
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util/cast"
	gconfig "github.com/snail007/gmc/util/config"
	"reflect"
)

var (
	groupMySQL   = NewMySQLDBGroup("default")
	groupSQLite3 = NewSQLite3DBGroup("default")
	cfg          *gconfig.Config
)

//InitFromFile parse foo.toml database configuration, `cfg` is Config object of foo.toml
func InitFromFile(cfgFile string) (err error) {
	cfg := gconfig.NewConfig()
	cfg.SetConfigFile(cfgFile)
	err = cfg.ReadInConfig()
	if err != nil {
		return err
	}
	return Init(cfg)
}

//RegistGroup parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 *gconfig.Config) (err error) {
	cfg = cfg0
	for k, v := range cfg.Sub("database").AllSettings() {
		if _, ok := v.([]interface{}); !ok {
			continue
		}
		for _, vv := range v.([]interface{}) {
			vvv := vv.(map[string]interface{})
			if !gcast.ToBool(vvv["enable"]) {
				continue
			}
			id := gcast.ToString(vvv["id"])
			if k == "mysql" {
				db := groupMySQL.DB(id)
				if db != nil {
					return
				}
				err = groupMySQL.Regist(id, MySQLDBConfig{
					Charset:                  gcast.ToString(vvv["charset"]),
					Collate:                  gcast.ToString(vvv["collate"]),
					Host:                     gcast.ToString(vvv["host"]),
					Port:                     gcast.ToInt(vvv["port"]),
					Database:                 gcast.ToString(vvv["database"]),
					Username:                 gcast.ToString(vvv["username"]),
					Password:                 gcast.ToString(vvv["password"]),
					TablePrefix:              gcast.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: gcast.ToString(vvv["prefix_sql_holder"]),
					Timeout:                  gcast.ToInt(vvv["timeout"]),
					ReadTimeout:              gcast.ToInt(vvv["readtimeout"]),
					WriteTimeout:             gcast.ToInt(vvv["writetimeout"]),
					SetMaxIdleConns:          gcast.ToInt(vvv["maxidle"]),
					SetMaxOpenConns:          gcast.ToInt(vvv["maxconns"]),
				})
				if err != nil {
					return
				}
			} else if k == "sqlite3" {
				db := groupSQLite3.DB(id)
				if db != nil {
					return
				}
				err = groupSQLite3.Regist(id, SQLite3DBConfig{
					Database:                 gcast.ToString(vvv["database"]),
					Password:                 gcast.ToString(vvv["password"]),
					TablePrefix:              gcast.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: gcast.ToString(vvv["prefix_sql_holder"]),
					SyncMode:                 gcast.ToInt(vvv["syncmode"]),
					OpenMode:                 gcast.ToString(vvv["openmode"]),
					CacheMode:                gcast.ToString(vvv["cachemode"]),
				})
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func DB(id ...string) gcore.Database {
	switch cfg.GetString("database.default") {
	case "mysql":
		return DBMySQL(id...)
	case "sqlite3":
		return DBSQLite3(id...)
	}
	return nil
}

//DBMySQL acquires a mysql db object associated the id, id default is : `default`
func DBMySQL(id ...string) *MySQLDB {
	return groupMySQL.DB(id...).(*MySQLDB)
}

//DBSQLite3 acquires a sqlite3 db object associated the id, id default is : `default`
func DBSQLite3(id ...string) *SQLite3DB {
	return groupSQLite3.DB(id...).(*SQLite3DB)
}

func isArray(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}
func isBool(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

package gmcdb

import (
	gmcconfig "github.com/snail007/gmc/config"
	"github.com/snail007/gmc/util/cast"
)

var (
	groupMySQL   = NewMySQLDBGroup("default")
	groupSQLite3 = NewSQLite3DBGroup("default")
	cfg          *gmcconfig.Config
)

//InitFromFile parse foo.toml database configuration, `cfg` is Config object of foo.toml
func InitFromFile(cfgFile string) (err error) {
	cfg := gmcconfig.New()
	cfg.SetConfigFile(cfgFile)
	err = cfg.ReadInConfig()
	if err != nil {
		return err
	}
	return Init(cfg)
}

//Init parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 *gmcconfig.Config) (err error) {
	cfg = cfg0
	for k, v := range cfg.Sub("database").AllSettings() {
		if _, ok := v.([]interface{}); !ok {
			continue
		}
		for _, vv := range v.([]interface{}) {
			vvv := vv.(map[string]interface{})
			if !cast.ToBool(vvv["enable"]) {
				continue
			}
			id := cast.ToString(vvv["id"])
			if k == "mysql" {
				db := groupMySQL.DB(id)
				if db != nil {
					return
				}
				err = groupMySQL.Regist(id, MySQLDBConfig{
					Charset:                  cast.ToString(vvv["charset"]),
					Collate:                  cast.ToString(vvv["collate"]),
					Host:                     cast.ToString(vvv["host"]),
					Port:                     cast.ToInt(vvv["port"]),
					Database:                 cast.ToString(vvv["database"]),
					Username:                 cast.ToString(vvv["username"]),
					Password:                 cast.ToString(vvv["password"]),
					TablePrefix:              cast.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: cast.ToString(vvv["prefix_sql_holder"]),
					Timeout:                  cast.ToInt(vvv["timeout"]),
					ReadTimeout:              cast.ToInt(vvv["readtimeout"]),
					WriteTimeout:             cast.ToInt(vvv["writetimeout"]),
					SetMaxIdleConns:          cast.ToInt(vvv["maxidle"]),
					SetMaxOpenConns:          cast.ToInt(vvv["maxconns"]),
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
					Database:                 cast.ToString(vvv["database"]),
					Password:                 cast.ToString(vvv["password"]),
					TablePrefix:              cast.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: cast.ToString(vvv["prefix_sql_holder"]),
					SyncMode:                 cast.ToInt(vvv["syncmode"]),
					OpenMode:                 cast.ToString(vvv["openmode"]),
					CacheMode:                cast.ToString(vvv["cachemode"]),
				})
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func DB(id ...string) Database {
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

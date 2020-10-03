package gmcdbhelper

import (
	gmcconfig "github.com/snail007/gmc/config"
	gmcmysql "github.com/snail007/gmc/db/mysql"
	gmcsqlite3 "github.com/snail007/gmc/db/sqlite3"
	"github.com/snail007/gmc/util/castutil"
)

var (
	groupMySQL   = *gmcmysql.NewDBGroup("default")
	groupSQLite3 = *gmcsqlite3.NewDBGroup("default")
	cfg        *gmcconfig.Config
)

//RegistGroup parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 *gmcconfig.Config) (err error) {
	cfg=cfg0
	for k, v := range cfg.Sub("database").AllSettings() {
		if _, ok := v.([]interface{}); !ok {
			continue
		}
		for _, vv := range v.([]interface{}) {
			vvv := vv.(map[string]interface{})
			if !castutil.ToBool(vvv["enable"]) {
				continue
			}
			id := castutil.ToString(vvv["id"])
			if k == "mysql" {
				db := groupMySQL.DB(id)
				if db != nil {
					return
				}
				err = groupMySQL.Regist(id, gmcmysql.DBConfig{
					Charset:                  castutil.ToString(vvv["charset"]),
					Collate:                  castutil.ToString(vvv["collate"]),
					Host:                     castutil.ToString(vvv["host"]),
					Port:                     castutil.ToInt(vvv["port"]),
					Database:                 castutil.ToString(vvv["database"]),
					Username:                 castutil.ToString(vvv["username"]),
					Password:                 castutil.ToString(vvv["password"]),
					TablePrefix:              castutil.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: castutil.ToString(vvv["prefix_sql_holder"]),
					Timeout:                  castutil.ToInt(vvv["timeout"]),
					ReadTimeout:              castutil.ToInt(vvv["readtimeout"]),
					WriteTimeout:             castutil.ToInt(vvv["writetimeout"]),
					SetMaxIdleConns:          castutil.ToInt(vvv["maxidle"]),
					SetMaxOpenConns:          castutil.ToInt(vvv["maxconns"]),
				})
				if err != nil {
					return
				}
			} else if k == "sqlite3" {
				db := groupSQLite3.DB(id)
				if db != nil {
					return
				}
				err = groupSQLite3.Regist(id, gmcsqlite3.DBConfig{
					Database:                 castutil.ToString(vvv["database"]),
					Password:                 castutil.ToString(vvv["password"]),
					TablePrefix:              castutil.ToString(vvv["prefix"]),
					TablePrefixSqlIdentifier: castutil.ToString(vvv["prefix_sql_holder"]),
					SyncMode:                 castutil.ToInt(vvv["syncmode"]),
					OpenMode:                 castutil.ToString(vvv["openmode"]),
					CacheMode:                castutil.ToString(vvv["cachemode"]),
				})
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func DB(id ...string) interface{} {
	switch cfg.GetString("database.default") {
	case "mysql":
		return DBMySQL(id...)
	case "sqlite3":
		return DBSQLite3(id...)
	}
	return nil
}

//DBMySQL acquires a mysql db object associated the id, id default is : `default`
func DBMySQL(id ...string) *gmcmysql.DB {
	return groupMySQL.DB(id...)
}

//DBSQLite3 acquires a sqlite3 db object associated the id, id default is : `default`
func DBSQLite3(id ...string) *gmcsqlite3.DB {
	return groupSQLite3.DB(id...)
}

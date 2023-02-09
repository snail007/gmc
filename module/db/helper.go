// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"reflect"
	"sort"
	"strings"

	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
)

var (
	groupMySQL   = NewMySQLDBGroup("default")
	groupSQLite3 = NewSQLite3DBGroup("default")
	cfg          gcore.Config
	defaultDB    string
)

type M map[string]interface{}

//InitFromFile parse foo.toml database configuration, `cfg` is Config object of foo.toml
func InitFromFile(cfgFile string) (err error) {
	cfg := gcore.ProviderConfig()()
	cfg.SetConfigFile(cfgFile)
	err = cfg.ReadInConfig()
	if err != nil {
		return err
	}
	return Init(cfg)
}

//Init parse app.toml database configuration, `cfg` is Config object of app.toml
func Init(cfg0 gcore.Config) (err error) {
	defaultDB = cfg0.GetString("database.default")
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
					TablePrefixSQLIdentifier: gcast.ToString(vvv["prefix_sql_holder"]),
					Timeout:                  gcast.ToInt(vvv["timeout"]),
					ReadTimeout:              gcast.ToInt(vvv["readtimeout"]),
					WriteTimeout:             gcast.ToInt(vvv["writetimeout"]),
					MaxIdleConns:             gcast.ToInt(vvv["maxidle"]),
					MaxOpenConns:             gcast.ToInt(vvv["maxconns"]),
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
					TablePrefix:              gcast.ToString(vvv["prefix"]),
					TablePrefixSQLIdentifier: gcast.ToString(vvv["prefix_sql_holder"]),
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
	switch defaultDB {
	case "mysql":
		return DBMySQL(id...)
	case "sqlite3":
		return DBSQLite3(id...)
	}
	return nil
}

//DBMySQL acquires a mysql db object associated the id, id default is : `default`
func DBMySQL(id ...string) *MySQLDB {
	// no mysql database enabled, just return nil
	if len(groupMySQL.dbGroup) == 0 {
		return nil
	}
	return groupMySQL.DB(id...).(*MySQLDB)
}

//DBSQLite3 acquires a sqlite3 db object associated the id, id default is : `default`
func DBSQLite3(id ...string) *SQLite3DB {
	// no sqlite3 database enabled, just return nil
	if len(groupSQLite3.dbGroup) == 0 {
		return nil
	}
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

func sortMap(data gmap.M, asc bool) []map[string]interface{} {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	if asc {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}
	var m []gmap.M
	for _, k := range keys {
		col := k
		if strings.Index(k, ":") > 0 {
			col = k[strings.Index(k, ":")+1:]
		}
		m = append(m, gmap.M{
			"col":   col,
			"value": data[k],
		})
	}
	return m
}

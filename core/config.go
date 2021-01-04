// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"io"
	"os"
	"strings"
	"time"
)

type CommonConfig interface {
	SetConfigFile(in string)
	SetEnvPrefix(in string)
	AllowEmptyEnv(allowEmptyEnv bool)
	ConfigFileUsed() string
	AddConfigPath(in string)
	SetTypeByDefaultValue(enable bool)
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetSizeInBytes(key string) uint
	BindEnv(input ...string) error
	IsSet(key string) bool
	AutomaticEnv()
	SetEnvKeyReplacer(r *strings.Replacer)
	RegisterAlias(alias string, key string)
	InConfig(key string) bool
	SetDefault(key string, value interface{})
	Set(key string, value interface{})
	ReadInConfig() error
	MergeInConfig() error
	ReadConfig(in io.Reader) error
	MergeConfig(in io.Reader) error
	MergeConfigMap(cfg map[string]interface{}) error
	WriteConfig() error
	SafeWriteConfig() error
	WriteConfigAs(filename string) error
	SafeWriteConfigAs(filename string) error
	AllKeys() []string
	AllSettings() map[string]interface{}
	SetConfigName(in string)
	SetConfigType(in string)
	SetConfigPermissions(perm os.FileMode)
	Debug()
}

type SubConfig interface {
	CommonConfig
}

type Config interface {
	CommonConfig
	Sub(key string) SubConfig
}

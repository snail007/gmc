// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gi18n

import (
	"bytes"
	"encoding/base64"
	gcore "github.com/snail007/gmc/core"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"
)

var (
	bindata = map[string][]byte{}
	I18N    = newI18n()
)

func SetBinData(data map[string]string) {
	bindata = map[string][]byte{}
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic("init i8n bin data fail, error: " + err.Error())
		}
		bindata[k] = b
	}
}

func Init(cfg gcore.Config) (err error) {
	if cfg.Sub("i18n") == nil {
		return
	}
	if len(bindata) > 0 {
		return initFromBinData(cfg)
	}
	return initFromDisk(cfg)
}

func initFromBinData(cfg gcore.Config) (err error) {
	fallbackLang := cfg.GetString("i18n.default")
	enabled := cfg.GetBool("i18n.enable")
	if !enabled {
		return
	}
	I18N = &I18n{
		langs:        map[string]map[string]string{},
		fallbackLang: fallbackLang,
	}
	for lang, v := range bindata {
		c := gcore.Providers.Config("")()
		err = c.ReadConfig(bytes.NewReader(v))
		if err != nil {
			return
		}
		if _, e := language.Parse(lang); e != nil {
			return gcore.Providers.Error("")().New(e)
		}
		data := map[string]string{}
		for _, k := range c.AllKeys() {
			data[k] = c.GetString(k)
		}
		I18N.Add(lang, data)
	}
	return
}

func initFromDisk(cfg gcore.Config) (err error) {
	dir := cfg.GetString("i18n.dir")
	fallbackLang := cfg.GetString("i18n.default")
	enalbed := cfg.GetBool("i18n.enable")
	if !enalbed {
		return
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	if err != nil {
		return
	}
	I18N = &I18n{
		langs:        map[string]map[string]string{},
		fallbackLang: fallbackLang,
	}
	for _, f := range files {
		c := gcore.Providers.Config("")()
		c.SetConfigFile(f)
		err = c.ReadInConfig()
		if err != nil {
			return
		}
		lang := filepath.Base(f)
		lang = strings.TrimSuffix(lang, filepath.Ext(lang))
		if _, e := language.Parse(lang); e != nil {
			return gcore.Providers.Error("")().New(e)
		}
		data := map[string]string{}
		for _, k := range c.AllKeys() {
			data[k] = c.GetString(k)
		}
		I18N.Add(lang, data)
	}
	return
}

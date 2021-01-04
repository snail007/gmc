// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"golang.org/x/text/language"
	"html/template"
	"net/http"
)

type I18n interface {
	Add(lang string, data map[string]string)
	Lang(lang string)
	String(lang string)
	Tr(lang, key string, defaultMessage ...string) string
	TrLangs(langs []string, key string, defaultMessage ...string) string
	TrV(lang, key string, defaultMessage ...string) template.HTML
	ParseAcceptLanguageStr(s string) (str string, err error)
	ParseAcceptLanguageStrT(s string) (t language.Tag, err error)
	ParseAcceptLanguage(r *http.Request) (strs []string, err error)
	ParseAcceptLanguageT(r *http.Request) (tags []language.Tag, err error)
	Languages() (languages []string, err error)
	LanguagesT() (languages []language.Tag, err error)
	MatchT(languages []language.Tag) (t language.Tag, err error)
	MatchAcceptLanguageT(r *http.Request) (tag language.Tag, err error)
	Clone(lang string) I18n
}

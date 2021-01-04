// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gi18n

import (
	gcore "github.com/snail007/gmc/core"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"strings"
)

type I18n struct {
	langs        map[string]map[string]string
	fallbackLang string
}

func newI18n() *I18n {
	return &I18n{
		langs: map[string]map[string]string{},
	}
}

func (this *I18n) Clone(lang string) gcore.I18n {
	return &I18n{
		langs:        this.langs,
		fallbackLang: lang,
	}
}
func (this *I18n) Add(lang string, data map[string]string) {
	this.langs[strings.ToLower(lang)] = data
}

func (this *I18n) Lang(lang string) {
	this.fallbackLang = strings.ToLower(lang)
}

func (this *I18n) String(lang string) {
	this.fallbackLang = strings.ToLower(lang)
}

func (this *I18n) Tr(lang, key string, defaultMessage ...string) string {
	if lang == "" {
		lang = this.fallbackLang
	}
	msg := key
	if len(defaultMessage) > 0 {
		msg = defaultMessage[0]
	}
	for _, k := range []string{lang, this.fallbackLang} {
		if v, ok := this.langs[k]; ok {
			if vv, ok := v[key]; ok {
				return vv
			}
		}
	}
	return msg
}

func (this *I18n) TrLangs(langs []string, key string, defaultMessage ...string) string {
	langs = append(langs, this.fallbackLang)
	msg := key
	if len(defaultMessage) > 0 {
		msg = defaultMessage[0]
	}
	for _, k := range langs {
		if v, ok := this.langs[k]; ok {
			if vv, ok := v[key]; ok {
				return vv
			}
		}
	}
	return msg
}

func (this *I18n) TrV(lang, key string, defaultMessage ...string) template.HTML {
	return template.HTML(this.Tr(lang, key, defaultMessage...))
}

func (this *I18n) ParseAcceptLanguageStr(s string) (str string, err error) {
	t, err := language.Parse(s)
	if err == nil {
		str = t.String()
		return
	}
	return
}

func (this *I18n) ParseAcceptLanguageStrT(s string) (t language.Tag, err error) {
	return language.Parse(s)
}

func (this *I18n) ParseAcceptLanguage(r *http.Request) (strs []string, err error) {
	t, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if len(t) > 0 && err == nil {
		for _, v := range t {
			strs = append(strs, v.String())
		}
		return
	}
	return
}
func (this *I18n) ParseAcceptLanguageT(r *http.Request) (tags []language.Tag, err error) {
	t, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if len(t) > 0 && err == nil {
		for _, v := range t {
			tags = append(tags, v)
		}
		return
	}
	return
}

func (this *I18n) Languages() (languages []string, err error) {
	for l := range this.langs {
		languages = append(languages, language.MustParse(l).String())
	}
	return
}

func (this *I18n) LanguagesT() (languages []language.Tag, err error) {
	for l := range this.langs {
		languages = append(languages, language.MustParse(l))
	}
	return
}
func (this *I18n) MatchT(languages []language.Tag) (t language.Tag, err error) {
	supportedLangs, err := this.LanguagesT()
	if err != nil {
		return
	}
	m := language.NewMatcher(supportedLangs)
	t, _, _ = m.Match(languages...)
	return
}

func (this *I18n) MatchAcceptLanguageT(r *http.Request) (tag language.Tag, err error) {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return
	}
	return this.MatchT(tags)
}

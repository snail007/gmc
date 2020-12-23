package gi18n

import (
	gcore "github.com/snail007/gmc/core"
	"strings"
)

type I18nTool struct {
	lang string
	tr   func(lang, key string, defaultMessage ...string) string
}

func NewI18nTool(tr func(lang, key string, defaultMessage ...string) string) *I18nTool {
	return &I18nTool{
		tr: tr,
	}
}

func (this *I18nTool) Lang(lang string) *I18nTool {
	this.lang = strings.ToLower(lang)
	return this
}

func (this *I18nTool) Tr(key string, defaultMessage ...string) string {
	return this.tr(this.lang, key, defaultMessage...)
}

type I18n struct {
	langs       map[string]map[string]string
	fallbackLang string
}

func New() *I18n {
	return &I18n{
		langs: map[string]map[string]string{},
	}
}

func (this *I18n) Add(lang string, data map[string]string) gcore.I18n {
	this.langs[strings.ToLower(lang)] = data
	return this
}

func (this *I18n) Lang(lang string) gcore.I18n  {
	this.fallbackLang = strings.ToLower(lang)
	return this
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
	langs=append(langs,this.fallbackLang)
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
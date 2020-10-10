package gmci18n

import (
	gmcconfig "github.com/snail007/gmc/config"
	gmcerr "github.com/snail007/gmc/error"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	i18n = New()
)

func Init(cfg *gmcconfig.Config) (err error) {
	dir := cfg.GetString("i18n.dir")
	defaultLang := cfg.GetString("i18n.default")
	enalbed := cfg.GetBool("i18n.enable")
	if !enalbed {
		return
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	if err != nil {
		return
	}
	i18n.Lang(defaultLang)
	for _, f := range files {
		c := gmcconfig.New()
		c.SetConfigFile(f)
		err = c.ReadInConfig()
		if err != nil {
			return
		}
		lang := filepath.Base(f)
		lang = strings.TrimSuffix(lang, filepath.Ext(lang))
		if _, e := language.Parse(lang); e != nil {
			return gmcerr.New(err)
		}
		data := map[string]string{}
		for _, k := range c.AllKeys() {
			data[k] = c.GetString(k)
		}
		i18n.Add(lang, data)
	}
	return
}

func Tr(lang, key string, defaultMessage ...string) string {
	return i18n.Tr(lang, key, defaultMessage...)
}

func TrV(lang, key string, defaultMessage ...string) template.HTML {
	return template.HTML(i18n.Tr(lang, key, defaultMessage...))
}

func TrLangs(langs []string, key string, defaultMessage ...string) string {
	return i18n.TrLangs(langs, key, defaultMessage...)
}

func ParseAcceptLanguageStr(s string)(str string,err error) {
	t,err:= language.Parse(s)
	if err==nil{
		str=t.String()
		return
	}
	return
}

func ParseAcceptLanguageStrT(s string)(t language.Tag,err error) {
	return language.Parse(s)
}

func ParseAcceptLanguage(r *http.Request)(strs []string,err error)  {
	t,_,err:= language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if len(t)>0&&err==nil{
		for _,v:=range t{
			strs=append(strs,v.String())
		}
		return
	}
	return
}
func ParseAcceptLanguageT(r *http.Request)(tags []language.Tag,err error)  {
	t,_,err:= language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if len(t)>0&&err==nil{
		for _,v:=range t{
			tags=append(tags,v)
		}
		return
	}
	return
}

func Languages() (languages []string,err error) {
	for l := range i18n.langs {
		languages = append(languages, language.MustParse(l).String())
	}
	return
}

func LanguagesT() (languages []language.Tag,err error) {
	for l := range i18n.langs {
		languages = append(languages, language.MustParse(l))
	}
	return
}

func MatchT (languages []language.Tag)(t language.Tag,err error)  {
	supportedLangs,err:=LanguagesT()
	if err!=nil{
		return
	}
	m:=language.NewMatcher(supportedLangs)
	t,_,_=m.Match(languages...)
	return
}

func MatchAcceptLanguageT(r *http.Request)(tag language.Tag,err error)  {
	tags,_,err:= language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err!=nil{
		return
	}
	return MatchT(tags)
}
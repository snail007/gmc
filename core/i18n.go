package gcore

type I18n interface {
	Add(lang string, data map[string]string)
	Lang(lang string)
	Tr(lang, key string, defaultMessage ...string) string
	TrLangs(langs []string, key string, defaultMessage ...string) string
}

package gcore

type I18n interface {
	Add(lang string, data map[string]string) I18n
	Lang(lang string) I18n
	Tr(lang, key string, defaultMessage ...string) string
	TrLangs(langs []string, key string, defaultMessage ...string) string
}

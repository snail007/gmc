package gi18n

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"net/http"
	"testing"
)

func TestLanguages(t *testing.T) {
	i18n := newI18n()
	i18n.Add("en", map[string]string{})
	i18n.Add("fr", map[string]string{})

	result, err := i18n.Languages()

	assert.NoError(t, err)
	assert.Equal(t, []string{"en", "fr"}, result)
}

func TestParseAcceptLanguageT(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "en-US")

	result, err := I18N.ParseAcceptLanguageT(r)

	assert.NoError(t, err)
	assert.Equal(t, []language.Tag{language.MustParse("en-US")}, result)
}

func TestTrLangs(t *testing.T) {
	i18n := newI18n()
	i18n.Add("en", map[string]string{"key": "Hello"})
	i18n.Add("fr", map[string]string{"key": "Bonjour"})

	result := i18n.TrLangs([]string{"fr", "en"}, "key", "Default message")

	assert.Equal(t, "Bonjour", result)
}

func TestTrV(t *testing.T) {
	i18n := newI18n()
	i18n.Add("en", map[string]string{"key": "<b>Hello</b>"})

	result := i18n.TrV("en", "key", "Default message")

	assert.Equal(t, "<b>Hello</b>", string(result))
}

func TestParseAcceptLanguageStrT(t *testing.T) {
	result, err := I18N.ParseAcceptLanguageStrT("en-US")

	assert.NoError(t, err)
	assert.Equal(t, language.MustParse("en-US"), result)
}

func TestClone(t *testing.T) {
	i18n := newI18n()
	i18n.Add("en", map[string]string{"key": "Hello"})

	clone := i18n.Clone("en")

	assert.Equal(t, i18n.langs, clone.(*I18n).langs)
	assert.Equal(t, "en", clone.(*I18n).fallbackLang)
}

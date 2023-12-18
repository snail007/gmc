package gtemplate

import (
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender(t *testing.T) {
	tpl := `a{{.user}}c`
	d, err := RenderString(tpl, gmap.M{"user": "b"})
	assert.Nil(t, err)
	assert.Equal(t, "abc", string(d))

	tpl = `a {{ .user}c`
	_, err = RenderString(tpl, gmap.M{"user": "b"})
	assert.NotNil(t, err)
}

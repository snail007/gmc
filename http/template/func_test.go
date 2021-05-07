package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestVal(t *testing.T) {
	assert := assert2.New(t)
	ctx := gcore.ProviderCtx()()
	ctx.SetConfig(gcore.ProviderConfig()())
	tpl, _ := NewTemplate(ctx, "tests/views")
	tpl.Delims("{{", "}}")
	tpl.Funcs(map[string]interface{}{
		"val": trimNoValue,
	})
	tpl.Parse()
	output, err := tpl.Execute("val/val", map[string]interface{}{})
	assert.Nil(err)
	assert.Empty(output)

	output, err = tpl.Execute("val/val1", map[string]string{
		"test": "test",
	})
	assert.Nil(err)
	assert.Equal("test", string(output))

	output, err = tpl.Execute("val/val2", map[string]interface{}{
		"map": map[string]string{"a": "a"},
	})
	assert.Nil(err)
	assert.Equal("a", string(output))

	output, err = tpl.Execute("val/val2", map[string]interface{}{
		"map": map[string]interface{}{"a": "a"},
	})
	assert.Nil(err)
	assert.Equal("a", string(output))

	output, err = tpl.Execute("val/val2", map[string]interface{}{
		"map": map[string]interface{}{"b": "b"},
	})
	assert.Nil(err)
	assert.Empty(output)
}

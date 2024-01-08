package gtemplate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender_Parse(t *testing.T) {
	// 创建一个 Render 实例
	renderer := NewRender()

	// 定义模板内容
	tplContent := `Hello, {{.Name}}!`

	// 定义模板数据
	tplData := map[string]interface{}{
		"Name": "World",
	}

	// 调用 Parse 方法解析模板
	result, err := renderer.Parse([]byte(tplContent), tplData)
	if err != nil {
		t.Fatalf("Error parsing template: %v", err)
	}

	// 验证结果是否符合预期
	expectedResult := "Hello, World!"
	if string(result) != expectedResult {
		t.Errorf("Expected result: %s, but got: %s", expectedResult, result)
	}

	// 测试 Delims 方法
	renderer.Delims("[[", "]]")
	result, err = renderer.Parse([]byte(`Hello, [[.Name]]!`), tplData)
	if err != nil {
		t.Fatalf("Error parsing template with custom delimiters: %v", err)
	}

	expectedResult = "Hello, World!"
	if string(result) != expectedResult {
		t.Errorf("Expected result with custom delimiters: %s, but got: %s", expectedResult, result)
	}
	renderer.Delims("{{", "}}")

	// 测试 AddFuncMap 方法
	funcMap := map[string]interface{}{
		"Double": func(s string) string {
			return s + s
		},
	}
	renderer.AddFuncMap(funcMap)
	result, err = renderer.Parse(`Doubled: {{Double .Name}}`, tplData)
	if err != nil {
		t.Fatalf("Error parsing template with added function: %v", err)
	}

	expectedResult = "Doubled: WorldWorld"
	if string(result) != expectedResult {
		t.Errorf("Expected result with added function: %s, but got: %s", expectedResult, result)
	}

	_, err = renderer.Parse(`Doubled: {{Double .Name}}`, tplData)
	assert.Nil(t, err)

	_, err = renderer.Parse(123, nil)
	assert.NotNil(t, err)

	_, err = renderer.Parse(`{{none .Name}}`, nil)
	assert.NotNil(t, err)

	_, err = renderer.Parse(`{{none .Name}`, nil)
	assert.NotNil(t, err)

	renderer.AddFuncMap(nil)
	assert.NotEmpty(t, renderer.funcMap)
}

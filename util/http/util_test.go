package ghttp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetContents(t *testing.T) {
	assert.Equal(t, "hello", string(GetContent(httpServerURL+"/hello", time.Second)))
	_, e := GetContentE("http://none/none", time.Second)
	assert.NotNil(t, e)
}

package gmcapp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	cfg := NewConfig()

	//run gmc app
	app := New(cfg)
	app.BeforeRun(func() error {
		return fmt.Errorf("stop")
	})
	err := app.Run()
	assert.Equal(err.Error(), "stop")
}

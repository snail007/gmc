package gtest

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNewProcess(t *testing.T) {
	t.Run("gtest.StartAndKill", func(t *testing.T) {
		assert := assert2.New(t)
		if RunProcess(t, func() {
			fmt.Println("abc")
			select {}
		}) {
			return
		}
		if InGMCT() {
			os.Setenv("GMCT_COVER_SHOW_KILLED", "cover_killed")
			os.Setenv("GMCT_COVER_VERBOSE", "true")
		}
		p := NewProcess(t)
		p.Start()
		time.Sleep(time.Second * 3)
		//assert.True(p.IsRunning())
		p.Kill()
		assert.Contains(p.Output(), "abc")
		if InGMCT() {
			assert.Contains(p.Output(), "cover_killed")
		}
	})
	t.Run("gtest.Wait", func(t *testing.T) {
		assert := assert2.New(t)
		if RunProcess(t, func() {
			fmt.Println("abc")
		}) {
			return
		}
		os.Setenv("GMCT_COVER_VERBOSE", "true")
		p := NewProcess(t)
		out,code,err:=p.Wait()
		assert.Nil(err)
		assert.Equal(0,code)
		assert.Contains(out, "abc")
	})
}

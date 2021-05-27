package gtest

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestStartAndKill(t *testing.T) {
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
	assert.Contains(p.Kill().Error(), "not run")
	p.Start()
	_, _, e := p.Wait()
	assert.Contains(e.Error(), "already started")
	e = p.Start()
	assert.Contains(e.Error(), "already started")
	time.Sleep(time.Second * 3)
	assert.True(p.IsRunning())
	p.Kill()
	assert.Contains(p.Output(), "abc")
	if InGMCT() {
		assert.Contains(p.Output(), "cover_killed")
	}
}
func TestWait(t *testing.T) {
	assert := assert2.New(t)
	if RunProcess(t, func() {
		fmt.Println("abc")
	}) {
		return
	}
	os.Setenv("GMCT_COVER_VERBOSE", "true")
	p := NewProcess(t)
	out, code, err := p.Wait()
	assert.Nil(err)
	assert.Equal(0, code)
	assert.Contains(out, "abc")
	assert.PanicsWithValue("abc", func() {
		panicErr("abc")
	})
}

func TestDebugRunProcess(t *testing.T) {
	k := execFlagPrefix + encodeTestName(t.Name())
	DebugRunProcess(t)
	assert2.Contains(t, os.Getenv(k), "true")
	os.Unsetenv(k)
}

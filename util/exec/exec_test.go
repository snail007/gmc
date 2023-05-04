package gexec

import (
	"bytes"
	_ "github.com/snail007/gmc"
	glog "github.com/snail007/gmc/module/log"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	//"which bash", 0, false, log
	output, errStr := NewCommand("which bash").Exec()
	assert.Empty(errStr)
	assert.Equal("/bin/bash\n", string(output))
}

func TestExec1(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	_, errStr := NewCommand("sleep 2").Timeout(time.Second).Exec()
	assert.NotEmpty(errStr)
}

func TestExec2(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	buf := &bytes.Buffer{}
	log := glog.New()
	log.SetOutput(buf)
	output, errStr := NewCommand("sleep 2").Log(log).Timeout(time.Second).Async(true).Exec()
	time.Sleep(time.Second * 4)
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "killed")
}

func TestExec4(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	buf := &bytes.Buffer{}
	log := glog.New()
	log.SetOutput(buf)
	output, errStr := NewCommand("ps -ef|grep").Log(log).Async(true).Timeout(time.Second).Exec()
	time.Sleep(time.Second * 2)
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "usage: grep")
}

func TestExec5(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	dir := "/tmp/exec_test_tmp"
	os.Mkdir(dir, 0755)
	output, err := NewCommand("pwd").WorkDir(dir).Exec()
	os.Remove(dir)
	assert.Empty(err)
	assert.Contains(output, "/exec_test_tmp")
}

func TestExec6(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	output, err := NewCommand("set|grep test_a|grep -v grep").Env(gmap.Mss{"test_a": "b"}).Exec()
	assert.Empty(err)
	assert.Contains(output, "test_a")
}

func TestExec7(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	output := ""
	var err error
	out, e := NewCommand("sleep 3; echo tst").Async(true).AsyncCallback(func(cmd *Command, msg string, e error) {
		output = msg
		err = e
	}).Exec()
	assert.Nil(e)
	assert.Empty(out)
	time.Sleep(time.Second * 4)
	assert.Nil(err)
	assert.Contains(output, "tst")
}

func TestExec8(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	buf := &bytes.Buffer{}
	log := glog.New()
	log.SetOutput(buf)
	output, errStr := NewCommand("ps aux").Log(log).Async(true).AsyncCallback(func(cmd *Command, output string, err error) {
		panic("test crash")
	}).Exec()
	time.Sleep(time.Second * 2)
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "test crash")
}

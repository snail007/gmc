package gexec

import (
	"bytes"
	_ "github.com/snail007/gmc"
	glog "github.com/snail007/gmc/module/log"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	//"which bash", 0, false, log
	output, errStr := NewCommand("which bash").Exec()
	assert.Empty(errStr)
	assert.Contains(output, "/bin/bash\n")
}

func TestExec0(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	output, errStr := NewCommand("echo  $1").Args("arg1").Exec()
	assert.Empty(errStr)
	assert.Contains(output, "arg1")
}

func TestExec00(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	buf := bytes.Buffer{}
	output, errStr := NewCommand("echo  $1").Args("arg1").Output(&buf).Exec()
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "arg1")
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
	log.SetOutput(glog.NewLoggerWriter(buf))
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
	log.SetOutput(glog.NewLoggerWriter(buf))
	output, errStr := NewCommand("ps -ef|grep").Log(log).Async(true).Timeout(time.Second).Exec()
	time.Sleep(time.Second * 2)
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "exit")
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
	log.SetOutput(glog.NewLoggerWriter(buf))
	output, errStr := NewCommand("ps aux").Log(log).Async(true).AsyncCallback(func(cmd *Command, output string, err error) {
		panic("test crash")
	}).Exec()
	time.Sleep(time.Second * 2)
	assert.Empty(errStr)
	assert.Empty(output)
	assert.Contains(buf.String(), "test crash")
}

func TestCommand_Cmd(t *testing.T) {
	c := NewCommand("echo a")
	_, err := c.Exec()
	assert.Nil(t, err)
	assert.NotNil(t, c.Cmd())
}

func TestCommand_Kill_1(t *testing.T) {
	c := NewCommand("sleep 100").Async(true)
	a := ""
	b := ""
	d := ""
	c.BeforeExec(func(command *Command, cmd *exec.Cmd) {
		a = "111"
	})
	c.AfterExec(func(command *Command, cmd *exec.Cmd, err error) {
		b = "222"
	})
	c.AfterExited(func(command *Command, cmd *exec.Cmd, err error) {
		d = "333"
	})
	c.Kill()
	c.Exec()
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "111", a)
	assert.Equal(t, "222", b)
	assert.Equal(t, "", d)
	assert.False(t, c.Cmd().ProcessState != nil)
	c.Kill()
	time.Sleep(time.Millisecond * 100)
	assert.True(t, c.Cmd().ProcessState != nil)
	c.Kill()
}

func TestCommand_Hook(t *testing.T) {
	c := NewCommand("sleep 1")
	a := ""
	b := ""
	d := ""
	c.BeforeExec(func(command *Command, cmd *exec.Cmd) {
		a = "111"
	})
	c.AfterExec(func(command *Command, cmd *exec.Cmd, err error) {
		b = "222"
	})
	c.AfterExited(func(command *Command, cmd *exec.Cmd, err error) {
		d = "333"
	})
	c.Exec()
	assert.Equal(t, "111", a)
	assert.Equal(t, "222", b)
	assert.Equal(t, "333", d)
}

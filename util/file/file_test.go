package gfile

import (
	"github.com/magiconair/properties/assert"
	gvalue "github.com/snail007/gmc/util/value"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestAbs(t *testing.T) {
	path := Abs("file.go")
	assert2.Contains(t, path, "file.go")
	assert2.Contains(t, path, "file.go")
}

func TestBaseName(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal("c.txt", BaseName("/a/b/c.txt"))
	assert.Equal("c", BaseName("/a/b/c"))
	assert.Equal("b", BaseName("/a/b/"))
}

func TestBytes(t *testing.T) {
	assert := assert2.New(t)
	file := "foo.txt"
	defer os.Remove(file)
	assert.Nil(WriteString(file, "abc", false))
	s := Bytes(file)
	assert.Equal(s, []byte("abc"))
}

func TestExists(t *testing.T) {
	assert2.True(t, Exists("file.go"))
	assert2.True(t, Exists("../file"))
	assert2.False(t, Exists("foo"))
}

func TestFileName(t *testing.T) {
	assert.Equal(t, FileName("/a/b/c.txt"), "c")
	assert.Equal(t, FileName("/a/b/c"), "c")
	assert.Equal(t, FileName("/a/b/"), "b")
}

func TestFileSize(t *testing.T) {
	f := "foo.txt"
	defer os.Remove(f)
	_t := time.Now()
	assert2.Nil(t, WriteString(f, "a", false))
	assert.Equal(t, _t.Unix(), gvalue.Must(FileMTime(f)).Time().Unix())
	size, err := FileSize(f)
	assert2.Nil(t, err)
	assert.Equal(t, size, int64(1))
	size, err = FileSize(f + "1")
	assert2.NotNil(t, err)
	assert.Equal(t, size, int64(0))
}

func TestHomeDir(t *testing.T) {
	assert2.NotEmpty(t, HomeDir())
}

func TestIsDir(t *testing.T) {
	assert2.True(t, IsDir("../file"))
	assert2.False(t, IsDir("file.go"))
}

func TestIsEmptyDir(t *testing.T) {
	dir := "dirfoo"
	assert := assert2.New(t)
	assert.Nil(os.Mkdir(dir, 0755))
	defer os.Remove(dir)
	assert.True(IsEmptyDir(dir))
	assert.False(IsEmptyDir("."))
}

func TestIsFile(t *testing.T) {
	assert2.True(t, IsFile("file.go"))
	assert2.False(t, IsFile("foo.go"))
	assert2.False(t, IsFile("../file"))
}

func TestIsLink(t *testing.T) {
	assert := assert2.New(t)
	dir := "dirfoo"
	linkdir := dir + "1"
	defer os.Remove(dir)
	defer os.Remove(linkdir)
	assert.Nil(os.Mkdir(dir, 0755))
	assert.Nil(os.Symlink(dir, linkdir))
	assert.False(IsLink(dir))
	assert.True(IsLink(linkdir))
}

func TestResolveLink(t *testing.T) {
	assert := assert2.New(t)
	dir := "dirfoo"
	linkdir := dir + "1"
	defer os.Remove(dir)
	defer os.Remove(linkdir)
	assert.Nil(os.Mkdir(dir, 0755))
	assert.Nil(os.Symlink(dir, linkdir))
	path, err := ResolveLink(dir)
	assert.Nil(err)
	assert.Contains(path, dir)
	path, err = ResolveLink(linkdir)
	assert.Nil(err)
	assert.Contains(path, dir)
	assert.NotContains(path, linkdir)
}

func TestReadAll(t *testing.T) {
	assert := assert2.New(t)
	file := "foo.txt"
	defer os.Remove(file)
	assert.Nil(WriteString(file, "abc", true))
	s := ReadAll(file)
	assert.Equal(s, "abc")
}

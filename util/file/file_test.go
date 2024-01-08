package gfile

import (
	"github.com/magiconair/properties/assert"
	gvalue "github.com/snail007/gmc/util/value"
	assert2 "github.com/stretchr/testify/assert"
	"io/ioutil"
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

func TestCopyFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "copy_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	sourceFilePath := tempDir + "/source.txt"
	destinationFilePath := tempDir + "/destination.txt"

	err = ioutil.WriteFile(sourceFilePath, []byte("Hello, world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Copy(sourceFilePath, destinationFilePath, true)
	if err != nil {
		t.Fatalf("Error copying file: %v", err)
	}

	destinationContent, err := ioutil.ReadFile(destinationFilePath)
	if err != nil {
		t.Fatalf("Error reading destination file: %v", err)
	}

	expectedContent := []byte("Hello, world!")
	if string(destinationContent) != string(expectedContent) {
		t.Fatalf("Copied file content does not match. Expected %q, got %q", expectedContent, destinationContent)
	}

	nonExistentDirPath := tempDir + "/nonexistent"
	nonExistentDestinationPath := nonExistentDirPath + "/destination.txt"

	err = Copy(sourceFilePath, nonExistentDestinationPath, false)
	if err == nil {
		t.Fatal("Expected error when destination directory does not exist, but got none")
	} else {
		t.Logf("Expected error: %v", err)
	}

	unwritableDirPath := "/unwritable"
	unwritableDestinationPath := unwritableDirPath + "/destination.txt"

	err = Copy(sourceFilePath, unwritableDestinationPath, true)
	if err == nil {
		t.Fatal("Expected error when destination directory is unwritable, but got none")
	} else {
		t.Logf("Expected error: %v", err)
	}

	nonExistentSourcePath := tempDir + "/nonexistent/source.txt"
	err = Copy(nonExistentSourcePath, destinationFilePath, true)
	if err == nil {
		t.Fatal("Expected error when source file does not exist, but got none")
	} else {
		t.Logf("Expected error: %v", err)
	}
}

func TestContent(t *testing.T) {
	// Create a temporary file for testing
	content := "This is a test file content."
	tmpFile, err := ioutil.TempFile("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to the temporary file
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Run the Content function with the temporary file
	result := Content(tmpFile.Name())

	// Check if the result matches the expected content
	if result != content {
		t.Errorf("Content(%s) = %s, want %s", tmpFile.Name(), result, content)
	}

	// Test case for a non-existing file
	nonExistentFile := "nonexistentfile.txt"
	result = Content(nonExistentFile)

	// Check if the result is an empty string for a non-existing file
	if result != "" {
		t.Errorf("Content(%s) = %s, want %s", nonExistentFile, result, "")
	}
}

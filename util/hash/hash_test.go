package ghash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func createTempFile(content string) (string, func(), error) {
	file, err := os.Create("testfile")
	if err != nil {
		return "", nil, err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return "", nil, err
	}
	return file.Name(), func() { os.Remove(file.Name()) }, nil
}

func TestMD5(t *testing.T) {
	input := "hello world"
	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte(input)))
	result := MD5(input)
	if result != expectedHash {
		t.Errorf("MD5 hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestMd5Bytes(t *testing.T) {
	input := []byte("hello world")
	expectedHash := fmt.Sprintf("%x", md5.Sum(input))
	result := Md5Bytes(input)
	if result != expectedHash {
		t.Errorf("Md5Bytes hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestMD5File(t *testing.T) {
	content := "test file content"
	filePath, cleanup, err := createTempFile(content)
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer cleanup()

	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte(content)))
	result, err := MD5File(filePath)
	if err != nil {
		t.Fatalf("Error calculating MD5 for file: %v", err)
	}

	if result != expectedHash {
		t.Errorf("MD5File hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA256(t *testing.T) {
	input := "hello world"
	expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
	result := SHA256(input)
	if result != expectedHash {
		t.Errorf("SHA256 hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA256File(t *testing.T) {
	content := "test file content"
	filePath, cleanup, err := createTempFile(content)
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer cleanup()

	expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	result, err := SHA256File(filePath)
	if err != nil {
		t.Fatalf("Error calculating SHA256 for file: %v", err)
	}

	if result != expectedHash {
		t.Errorf("SHA256File hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA1(t *testing.T) {
	input := "hello world"
	expectedHash := fmt.Sprintf("%x", sha1.Sum([]byte(input)))
	result := SHA1(input)
	if result != expectedHash {
		t.Errorf("SHA1 hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA1File(t *testing.T) {
	content := "test file content"
	filePath, cleanup, err := createTempFile(content)
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer cleanup()

	expectedHash := fmt.Sprintf("%x", sha1.Sum([]byte(content)))
	result, err := SHA1sumFile(filePath)
	if err != nil {
		t.Fatalf("Error calculating SHA1 for file: %v", err)
	}

	if result != expectedHash {
		t.Errorf("SHA1sumFile hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestBlake2s256(t *testing.T) {
	input := "hello world"
	expectedHash := Blake2s256(input)
	result := Blake2s256(input)
	if result != expectedHash {
		t.Errorf("Blake2s256 hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestBlake2s256File(t *testing.T) {
	content := "test file content"
	filePath, cleanup, err := createTempFile(content)
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer cleanup()

	expectedHash := Blake2s256(content)
	result, err := Blake2s256File(filePath)
	if err != nil {
		t.Fatalf("Error calculating Blake2s256 for file: %v", err)
	}

	if result != expectedHash {
		t.Errorf("Blake2s256File hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestCRC32(t *testing.T) {
	input := "hello world"
	expectedHash := CRC32(input)
	result := CRC32(input)
	if result != expectedHash {
		t.Errorf("CRC32 hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestCRC32File(t *testing.T) {
	content := "test file content"
	filePath, cleanup, err := createTempFile(content)
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer cleanup()

	expectedHash := CRC32(content)
	result, err := CRC32File(filePath)
	if err != nil {
		t.Fatalf("Error calculating CRC32 for file: %v", err)
	}

	if result != expectedHash {
		t.Errorf("CRC32File hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA256Bytes(t *testing.T) {
	input := []byte("hello world")
	expectedHash := fmt.Sprintf("%x", sha256.Sum256(input))
	result := SHA256Bytes(input)
	if result != expectedHash {
		t.Errorf("SHA256Bytes hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA1Bytes(t *testing.T) {
	input := []byte("hello world")
	expectedHash := fmt.Sprintf("%x", sha1.Sum(input))
	result := SHA1Bytes(input)
	if result != expectedHash {
		t.Errorf("SHA1Bytes hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestBlake2s256Bytes(t *testing.T) {
	input := []byte("hello world")
	expectedHash := Blake2s256Bytes(input)
	result := Blake2s256Bytes(input)
	if result != expectedHash {
		t.Errorf("Blake2s256Bytes hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestCRC32Bytes(t *testing.T) {
	input := []byte("hello world")
	expectedHash := CRC32Bytes(input)
	result := CRC32Bytes(input)
	if result != expectedHash {
		t.Errorf("CRC32Bytes hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestMD5Reader(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte(input)))
	result, err := MD5Reader(newReaderCloser(reader))
	if err != nil {
		t.Fatalf("Error calculating MD5 for reader: %v", err)
	}

	if result != expectedHash {
		t.Errorf("MD5Reader hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA256Reader(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	expectedHash := SHA256(input)
	result, err := SHA256Reader(newReaderCloser(reader))
	if err != nil {
		t.Fatalf("Error calculating SHA256 for reader: %v", err)
	}

	if result != expectedHash {
		t.Errorf("SHA256Reader hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestSHA1Reader(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	expectedHash := SHA1(input)
	result, err := SHA1Reader(newReaderCloser(reader))
	if err != nil {
		t.Fatalf("Error calculating SHA1 for reader: %v", err)
	}

	if result != expectedHash {
		t.Errorf("SHA1Reader hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestBlake2s256Reader(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	expectedHash := Blake2s256(input)
	result, err := Blake2s256Reader(newReaderCloser(reader))
	if err != nil {
		t.Fatalf("Error calculating Blake2s256 for reader: %v", err)
	}

	if result != expectedHash {
		t.Errorf("Blake2s256Reader hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestCRC32Reader(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	expectedHash := CRC32(input)
	result, err := CRC32Reader(newReaderCloser(reader))
	if err != nil {
		t.Fatalf("Error calculating CRC32 for reader: %v", err)
	}

	if result != expectedHash {
		t.Errorf("CRC32Reader hash mismatch. Expected: %s, Got: %s", expectedHash, result)
	}
}

func TestMd5Bytes_CloseError(t *testing.T) {
	t.Parallel()
	// Simulate a reader with an error when closing
	reader := &mockReadCloser{Reader: strings.NewReader("hello world"),
		Error: fmt.Errorf("error1")}
	_, err := MD5Reader(reader)
	if err == nil || err.Error() != "error1" {
		t.Errorf("Expected close error, but got: %v", err)
	}

	reader = &mockReadCloser{Reader: strings.NewReader("hello world"),
		sleep: time.Second * 33,
		Error: fmt.Errorf("error2")}
	_, err = MD5Reader(reader)
	if err == nil || err.Error() != "error2" {
		t.Errorf("Expected close error, but got: %v", err)
	}
}

func TestMD5File_ErrorOpeningFile(t *testing.T) {
	_, err := MD5File("nonexistentfile.txt")
	if err == nil {
		t.Error("Expected error opening file, but got nil")
	}
}

type mockReadCloser struct {
	io.Reader
	Error error
	sleep time.Duration
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	if m.sleep > 0 {
		time.Sleep(m.sleep)
	}
	return 0, m.Error
}

func (m *mockReadCloser) Close() error {
	return nil
}

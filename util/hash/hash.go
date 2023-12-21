package ghash

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	gsync "github.com/snail007/gmc/util/sync"
	"golang.org/x/crypto/blake2s"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const bufferSize = 65536

func MD5(str string) string {
	return Md5Bytes([]byte(str))
}

func Md5Bytes(d []byte) string {
	return fmt.Sprintf("%x", md5.Sum(d))
}

func MD5File(filename string) (string, error) {
	return sumFile(md5.New(), filename)
}

func MD5Reader(reader io.ReadCloser) (string, error) {
	return sumReader(md5.New(), reader)
}

func SHA256(str string) string {
	s, _ := sumReader(sha256.New(), newStringCloser(str))
	return s
}

func SHA256Bytes(b []byte) string {
	s, _ := sumReader(sha256.New(), newReaderCloser(bytes.NewReader(b)))
	return s
}

func SHA256File(filename string) (string, error) {
	return sumFile(sha256.New(), filename)
}

func SHA256Reader(reader io.ReadCloser) (string, error) {
	return sumReader(sha256.New(), reader)
}

func SHA1(str string) string {
	s, _ := sumReader(sha1.New(), newStringCloser(str))
	return s
}

func SHA1Bytes(b []byte) string {
	s, _ := sumReader(sha1.New(), newReaderCloser(bytes.NewReader(b)))
	return s
}

func SHA1sumFile(filename string) (string, error) {
	return sumFile(sha1.New(), filename)
}

func SHA1Reader(reader io.ReadCloser) (string, error) {
	return sumReader(sha1.New(), reader)
}

func Blake2s256(str string) string {
	hash, _ := blake2s.New256([]byte{})
	s, _ := sumReader(hash, newStringCloser(str))
	return s
}

func Blake2s256Bytes(b []byte) string {
	hash, _ := blake2s.New256([]byte{})
	s, _ := sumReader(hash, newReaderCloser(bytes.NewReader(b)))
	return s
}

func Blake2s256File(filename string) (string, error) {
	hash, _ := blake2s.New256([]byte{})
	return sumFile(hash, filename)
}

func Blake2s256Reader(reader io.ReadCloser) (string, error) {
	hash, _ := blake2s.New256([]byte{})
	return sumReader(hash, reader)
}

func CRC32(str string) string {
	s, _ := CRC32Reader(newStringCloser(str))
	return s
}

func CRC32Bytes(b []byte) string {
	s, _ := CRC32Reader(newReaderCloser(bytes.NewReader(b)))
	return s
}

func CRC32File(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil || info.IsDir() {
		return "", err
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	return CRC32Reader(bufio.NewReader(file))
}

func CRC32Reader(reader io.Reader) (string, error) {
	table := crc32.MakeTable(crc32.IEEE)
	checkSum := crc32.Checksum([]byte(""), table)
	buf := make([]byte, bufferSize)
	for {
		switch n, err := reader.Read(buf); err {
		case nil:
			checkSum = crc32.Update(checkSum, table, buf[:n])
		case io.EOF:
			return fmt.Sprintf("%08x", checkSum), nil
		default:
			return "", err
		}
	}
}

// sumReader calculates the hash based on a provided hash provider
func sumReader(hashAlgorithm hash.Hash, reader io.ReadCloser) (string, error) {
	buf := make([]byte, bufferSize)
	t := time.AfterFunc(time.Second*30, func() {
		reader.Close()
	})
	defer t.Stop()
	for {
		t.Reset(time.Second * 30)
		switch n, err := reader.Read(buf); err {
		case nil:
			hashAlgorithm.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hashAlgorithm.Sum(nil)), nil
		default:
			return "", err
		}
	}
}

type readerCloser struct {
	io.Reader
}

func newReaderCloser(r io.Reader) *readerCloser {
	return &readerCloser{Reader: r}
}

func (s *readerCloser) Close() error {
	return nil
}

func newStringCloser(str string) io.ReadCloser {
	return newReaderCloser(strings.NewReader(str))
}

// sumFile calculates the hash based on a provided hash provider
func sumFile(hashAlgorithm hash.Hash, filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil || info.IsDir() {
		return "", err
	}
	var file *os.File
	var err error
	g := sync.WaitGroup{}
	g.Add(1)
	go func() {
		defer g.Done()
		file, err = os.Open(filename)
	}()
	select {
	case <-gsync.Wait(&g):
	case <-time.After(time.Second * 3):
		return "", fmt.Errorf("open file timeout, %s", filename)
	}
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	return sumReader(hashAlgorithm, newBufReaderCloser(file))
}

type bufReaderCloser struct {
	*bufio.Reader
	f *os.File
}

func (s *bufReaderCloser) Close() error {
	if s.f != nil {
		return s.f.Close()
	}
	return nil
}

func newBufReaderCloser(f *os.File) *bufReaderCloser {
	return &bufReaderCloser{
		f:      f,
		Reader: bufio.NewReaderSize(f, 1024*8),
	}
}

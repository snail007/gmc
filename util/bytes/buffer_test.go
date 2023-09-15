package gbytes

import (
	"bytes"
	gcast "github.com/snail007/gmc/util/cast"
	gloop "github.com/snail007/gmc/util/loop"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestCircularBuffer_Write(t *testing.T) {
	buf := NewCircularBuffer(5)

	// Write data to the buffer
	n, err := buf.Write([]byte("12345"))
	if n != 5 || err != nil {
		t.Errorf("Expected n=5 and err=nil, got n=%d, err=%v", n, err)
	}
	assert.Equal(t, "12345", string(buf.Bytes()))

	r1 := buf.NewHistoryReader()
	readData := make([]byte, 3)
	n, err = r1.Read(readData)
	if n != 3 || err != nil {
		t.Errorf("Expected n=3 and err=nil, got n=%d, err=%v", n, err)
	}
	if string(readData) != "123" {
		t.Errorf("Expected '123', got '%s'", string(readData))
	}

	// Write more data to trigger overflow
	n, err = buf.Write([]byte("67890"))
	if n != 5 || err != nil {
		t.Errorf("Expected n=5 and err=nil, got n=%d, err=%v", n, err)
	}

	// Read from a reader
	reader := buf.NewReader()
	readData = make([]byte, 3)
	time.AfterFunc(time.Second, func() {
		buf.Write([]byte("000"))
	})
	n, err = reader.Read(readData)
	if n != 3 || err != nil {
		t.Errorf("Expected n=3 and err=nil, got n=%d, err=%v", n, err)
	}
	if string(readData) != "000" {
		t.Errorf("Expected '000', got '%s'", string(readData))
	}

	reader = buf.NewReader()
	readData = make([]byte, 6)
	time.AfterFunc(time.Second, func() {
		buf.Write([]byte("12345"))
	})
	n, err = reader.Read(readData)
	if n != 5 || err != nil {
		t.Errorf("Expected n=5 and err=nil, got n=%d, err=%v", n, err)
	}
	if string(readData[:5]) != "12345" {
		t.Errorf("Expected '12345', got '%s'", string(readData))
	}

	// Close the reader and write more data
	reader.Close()
	n, err = buf.Write([]byte("ABC"))
	if n != 3 || err != nil {
		t.Errorf("Expected n=3 and err=nil, got n=%d, err=%v", n, err)
	}

	// Try to write to a closed buffer
	reader = buf.NewReader()
	buf.Close()
	n, err = buf.Write([]byte("XYZ"))
	if n != 0 || err != io.ErrClosedPipe {
		t.Errorf("Expected n=0 and err=io.ErrClosedPipe, got n=%d, err=%v", n, err)
	}

	// Try to read from a closed buffer
	readData = make([]byte, 6)
	n, err = reader.Read(readData)
	if n != 0 || err == nil {
		t.Errorf("Expected n=0 and err!=nil, got n=%d, err=%v", n, err)
	}

	// Try to get reader from a closed buffer
	reader = buf.NewReader()
	readData = make([]byte, 6)
	n, err = reader.Read(readData)
	if n != 0 || err == nil {
		t.Errorf("Expected n=0 and err!=nil, got n=%d, err=%v", n, err)
	}
}

func TestCircularReader_Read(t *testing.T) {
	buf := NewCircularBuffer(5)
	reader := buf.NewReader()

	// Read from an empty reader
	readData := make([]byte, 3)
	time.AfterFunc(time.Second, func() {
		_ = reader.Close()
	})
	n, err := reader.Read(readData)
	if n != 0 || err == nil {
		t.Errorf("Expected n=0 and non-nil err, got n=%d, err=%v", n, err)
	}

	reader = buf.NewReader()

	// Write data to the buffer
	time.AfterFunc(time.Second, func() {
		buf.Write([]byte("12345"))
	})
	// Read from the reader
	n, err = reader.Read(readData)
	if n != 3 || err != nil {
		t.Errorf("Expected n=3 and err=nil, got n=%d, err=%v", n, err)
	}
	if string(readData) != "123" {
		t.Errorf("Expected '123', got '%s'", string(readData))
	}

	// Close the reader and try to read
	reader.Close()
	n, err = reader.Read(readData)
	if n != 0 || err != io.ErrClosedPipe {
		t.Errorf("Expected n=0 and err=io.ErrClosedPipe, got n=%d, err=%v", n, err)
	}
}

func TestCircularBuffer_Close(t *testing.T) {
	buf := NewCircularBuffer(5)
	reader := buf.NewReader()

	// Close the buffer while a reader is open
	err := buf.Close()
	if err != nil {
		t.Errorf("Expected err=nil, got err=%v", err)
	}

	// Try to read from the reader after buffer is closed
	readData := make([]byte, 3)
	n, err := reader.Read(readData)
	if n != 0 || err != io.ErrClosedPipe {
		t.Errorf("Expected n=0 and err=io.ErrClosedPipe, got n=%d, err=%v", n, err)
	}

	// Close the reader and try to close the buffer again
	reader.Close()
	err = buf.Close()
	if err != nil {
		t.Errorf("Expected err=nil, got err=%v", err)
	}
}

func TestCircularReader_Close(t *testing.T) {
	buf := NewCircularBuffer(5)
	reader := buf.NewReader()

	// Close the reader
	err := reader.Close()
	if err != nil {
		t.Errorf("Expected err=nil, got err=%v", err)
	}

	// Try to read from the closed reader
	readData := make([]byte, 3)
	n, err := reader.Read(readData)
	if n != 0 || err != io.ErrClosedPipe {
		t.Errorf("Expected n=0 and err=io.ErrClosedPipe, got n=%d, err=%v", n, err)
	}

	// Close the reader again (should not produce an error)
	err = reader.Close()
	if err != nil {
		t.Errorf("Expected err=nil, got err=%v", err)
	}
}

func TestCircularBuffer(t *testing.T) {

	// 1. 写入数据到CircularBuffer，包括正常情况和溢出情况。
	t.Run("Write", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)

		data1 := []byte("1234567")
		data2 := []byte("abcdef")

		n, err := cb.Write(data1)
		if err != nil {
			t.Errorf("Error writing data1: %v", err)
		}
		if n != len(data1) {
			t.Errorf("Expected to write %d bytes, but wrote %d bytes", len(data1), n)
		}

		n, err = cb.Write(data2)
		if err != nil {
			t.Errorf("Error writing data2: %v", err)
		}
		if n != len(data2) {
			t.Errorf("Expected to write %d bytes, but wrote %d bytes", len(data2), n)
		}

		expectedData := append(data1, data2...)
		if !bytes.Equal(cb.data, expectedData[len(expectedData)-bufferSize:]) {
			t.Errorf("CircularBuffer data not as expected: %v", cb.data)
		}
	})

	// 2. 创建读取器(CircularReader).
	t.Run("NewReader", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)
		reader := cb.NewReader()
		if reader == nil {
			t.Error("Failed to create CircularReader")
		}
	})

	// 3. 从读取器读取数据.
	t.Run("Read", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)

		reader := cb.NewReader()
		data := make([]byte, bufferSize)
		time.AfterFunc(time.Second, func() {
			reader.Close()
		})
		n, err := reader.Read(data)

		if err == nil {
			t.Errorf("Error reading data: %v", err)
		}
		if n != 0 {
			t.Errorf("Expected to read %d bytes, but read %d bytes", bufferSize, n)
		}
	})

	// 4. 关闭读取器和CircularBuffer.
	t.Run("Close", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)
		reader := cb.NewReader()
		err := reader.Close()
		if err != nil {
			t.Errorf("Error closing reader: %v", err)
		}
		go func() {
			time.Sleep(time.Second)
			cb.Close()
		}()
		reader = cb.NewReader()
		cb.SetReaderDeadline(reader, time.Now().Add(time.Minute))
		p := make([]byte, 1)
		_, err = reader.Read(p)
		assert.NotNil(t, err)
		err = cb.Close()
		if err != nil {
			t.Errorf("Error closing CircularBuffer: %v", err)
		}
	})

	t.Run("Close1", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)
		time.Sleep(time.Second)
		cb.Write([]byte("123"))
		cb.Reset()
		cb.Write([]byte("456"))
		assert.Equal(t, "456", string(cb.Bytes()))
	})
	t.Run("Close2", func(t *testing.T) {
		bufferSize := 10
		cb := NewCircularBuffer(bufferSize)
		reader := cb.NewReader()
		err := reader.Close()
		if err != nil {
			t.Errorf("Error closing reader: %v", err)
		}
		go func() {
			time.Sleep(time.Second)
			cb.Reset()
		}()
		reader = cb.NewReader()
		cb.SetReaderDeadline(reader, time.Now().Add(time.Minute))
		p := make([]byte, 1)
		_, err = reader.Read(p)
		assert.NotNil(t, err)
		err = cb.Close()
		if err != nil {
			t.Errorf("Error closing CircularBuffer: %v", err)
		}
	})
}

func TestCircularBuffer_Reset(t *testing.T) {
	// 创建一个CircularBuffer
	buffer := NewCircularBuffer(10)

	// 写入一些数据
	data := []byte("123456789")
	_, err := buffer.Write(data)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	reader := buffer.NewReader()
	readData := make([]byte, 3)
	buffer.ResetReader(reader)
	n, err := reader.Read(readData)
	if n != 3 || err != nil {
		t.Errorf("Expected n=3 and err=nil, got n=%d, err=%v", n, err)
	}
	if string(readData) != "123" {
		t.Errorf("Expected '123', got '%s'", string(readData))
	}
	buffer.Reset()
	// 重置后应该可以写入新数据
	_, err = buffer.Write([]byte("New Data"))
	if err != nil {
		t.Fatalf("Write error after reset: %v", err)
	}

	// 关闭CircularBuffer
	_ = buffer.Close()

	// 尝试在已关闭的情况下重置
	buffer.Reset()
}

func TestCircularReader_0(t *testing.T) {
	buf := NewCircularBuffer(5)
	reader := buf.NewReader()
	// Try to read from empty buffer
	buf.SetReaderDeadline(reader, time.Now().Add(time.Second))
	readData := make([]byte, 3)
	n, err := reader.Read(readData)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "exceeded")
	reader.Read(readData)
	buf.Write([]byte("123"))
	buf.SetReaderDeadline(reader, time.Time{})
	n, err = reader.Read(readData)
	assert.Equal(t, 3, n)
	assert.Equal(t, "123", string(readData))
	assert.Nil(t, err)
}

func TestCircularReader_1(t *testing.T) {
	buf := NewCircularBuffer(10)
	buf.Reset()
	go func() {
		gloop.For(20, func(loopIndex int) {
			time.Sleep(time.Millisecond)
			buf.Write([]byte(gcast.ToString(loopIndex + 1)))
		})
	}()
	g := sync.WaitGroup{}
	g.Add(1)
	go func() {
		defer g.Done()
		time.Sleep(time.Millisecond * 10)
		reader2 := buf.NewReader()
		buf.SetReaderDeadline(reader2, time.Now().Add(time.Second))
		b, err := ioutil.ReadAll(reader2)
		assert.NotNil(t, err)
		assert.Less(t, len(b), 31)
		assert.Greater(t, len(b), 10)
	}()
	reader := buf.NewHistoryReader()
	buf.SetReaderDeadline(reader, time.Now().Add(time.Second))
	b, err := ioutil.ReadAll(reader)
	assert.NotNil(t, err)
	assert.Equal(t, "1234567891011121314151617181920", string(b))
	g.Wait()
}

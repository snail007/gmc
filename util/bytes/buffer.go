package gbytes

import (
	"fmt"
	gerror "github.com/snail007/gmc/module/error"
	gmap "github.com/snail007/gmc/util/map"
	grand "github.com/snail007/gmc/util/rand"
	"io"
	"sync"
)

type CircularBuffer struct {
	data      []byte
	size      int
	isOpen    bool
	readers   []*CircularReader
	readersMu sync.Mutex
	waitQueue *gmap.Map
}

type CircularReader struct {
	buffer *CircularBuffer
	start  int
	closed bool
	waitCh chan bool
}

func NewCircularBuffer(size int) *CircularBuffer {
	b := &CircularBuffer{
		data:      make([]byte, 0, size),
		size:      size,
		isOpen:    true,
		waitQueue: gmap.New(),
	}
	return b
}
func (b *CircularBuffer) notify() {
	b.waitQueue.CloneAndClear().RangeFast(func(_, v interface{}) bool {
		b.closeCh(v.(chan bool))
		return true
	})
}
func (b *CircularBuffer) closeCh(ch chan bool) {
	gerror.Try(func() {
		close(ch)
	})
}
func (b *CircularBuffer) Write(p []byte) (n int, err error) {
	if !b.isOpen {
		return 0, io.ErrClosedPipe
	}

	b.data = append(b.data, p...)
	overflowCnt := 0
	if len(b.data) > b.size {
		overflowCnt = len(b.data) - b.size
		b.data = b.data[len(b.data)-b.size:]
	}

	b.readersMu.Lock()
	for _, r := range b.readers {
		r.start -= overflowCnt
		if r.start < 0 {
			r.start = 0
		}
	}
	b.readersMu.Unlock()
	b.notify()
	return len(p), nil
}

func (b *CircularBuffer) NewReader() io.ReadCloser {
	b.readersMu.Lock()
	defer b.readersMu.Unlock()

	if !b.isOpen {
		return &CircularReader{
			closed: true,
		}
	}

	for _, r := range b.readers {
		_ = r.Close()
	}
	b.readers = nil

	r := &CircularReader{buffer: b, start: 0}
	b.readers = append(b.readers, r)
	return r
}

func (r *CircularReader) Read(p []byte) (n int, err error) {
RETRY:
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	buffer := r.buffer
	if !buffer.isOpen {
		return 0, io.ErrClosedPipe
	}
	if len(buffer.data) == 0 {
		<-r.wait()
		goto RETRY
	}
	bufLen := len(buffer.data)
	start := r.start
	end := r.start + cap(p)
	if end >= bufLen {
		end = bufLen
	}
	n = copy(p, buffer.data[start:end])
	return
}

func (r *CircularReader) wait() chan bool {
	ch := make(chan bool)
	k := grand.String(8) + fmt.Sprintf("%p", ch)
	r.buffer.waitQueue.Store(k, ch)
	r.closeCh()
	r.waitCh = ch
	return ch
}

func (r *CircularReader) closeCh() {
	r.buffer.closeCh(r.waitCh)
}

func (r *CircularReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	r.closeCh()
	return nil
}

func (b *CircularBuffer) Close() error {
	b.isOpen = false
	b.readersMu.Lock()
	defer b.readersMu.Unlock()
	for _, r := range b.readers {
		_ = r.Close()
	}
	return nil
}

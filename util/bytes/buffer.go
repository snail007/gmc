package gbytes

import (
	"fmt"
	gerror "github.com/snail007/gmc/module/error"
	gcond "github.com/snail007/gmc/util/cond"
	gmap "github.com/snail007/gmc/util/map"
	grand "github.com/snail007/gmc/util/rand"
	"io"
	"sync"
	"time"
)

type CircularBuffer struct {
	data      []byte
	size      int
	isOpen    bool
	readers   []*CircularReader
	readersMu sync.Mutex
	waitQueue *gmap.Map
}

func NewCircularBuffer(size int) *CircularBuffer {
	b := &CircularBuffer{
		data:      []byte{},
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

func (b *CircularBuffer) Reset() {
	b.readersMu.Lock()
	defer b.readersMu.Unlock()
	for _, r := range b.readers {
		_ = r.Close()
	}
	b.waitQueue.RangeFast(func(_, v interface{}) bool {
		b.closeCh(v.(chan bool))
		return true
	})
	b.readers = []*CircularReader{}
	b.data = []byte{}
	b.waitQueue.Clear()
	b.waitQueue.GC()
}

func (b *CircularBuffer) ResetReader(r io.ReadCloser) {
	if v, ok := r.(*CircularReader); !ok {
		return
	} else {
		v.start = 0
	}
}

func (b *CircularBuffer) SetReaderDeadline(r io.ReadCloser, deadline time.Time) {
	if v, ok := r.(*CircularReader); ok {
		v.deadline = deadline
	}
}

func (b *CircularBuffer) Bytes() []byte {
	b.readersMu.Lock()
	defer b.readersMu.Unlock()
	if len(b.data) == 0 {
		return nil
	}
	bs := make([]byte, len(b.data))
	copy(bs, b.data)
	return bs
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

func (b *CircularBuffer) NewHistoryReader() io.ReadCloser {
	return b.newReader(false)
}

func (b *CircularBuffer) NewReader() io.ReadCloser {
	return b.newReader(true)
}

func (b *CircularBuffer) newReader(isCurrent bool) io.ReadCloser {
	b.readersMu.Lock()
	defer b.readersMu.Unlock()

	if !b.isOpen {
		return &CircularReader{
			closed: true,
		}
	}
	var readers []*CircularReader
	for _, r := range b.readers {
		if !r.closed {
			readers = append(readers, r)
		}
	}
	start := gcond.Cond(isCurrent, len(b.data)-1, 0).(int)
	r := &CircularReader{buffer: b, start: start}
	if r.start < 0 {
		r.start = 0
	}
	b.readers = append(b.readers, r)
	return r
}

func (b *CircularBuffer) Close() error {
	b.isOpen = false
	b.readersMu.Lock()
	defer b.readersMu.Unlock()
	for _, r := range b.readers {
		_ = r.Close()
	}
	b.waitQueue.RangeFast(func(_, v interface{}) bool {
		b.closeCh(v.(chan bool))
		return true
	})
	return nil
}

type CircularReader struct {
	buffer   *CircularBuffer
	start    int
	closed   bool
	waitCh   chan bool
	deadline time.Time
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
	if len(buffer.data) == 0 || r.start >= len(r.buffer.data)-1 {
		if e := r.wait(); e != nil {
			return 0, e
		}
		goto RETRY
	}
	bufLen := len(buffer.data)
	start := r.start
	end := r.start + cap(p)
	if end >= bufLen {
		end = bufLen
	}
	n = copy(p, buffer.data[start:end])
	r.start += n
	return
}

func (r *CircularReader) wait() error {
	ch := make(chan bool)
	k := grand.String(8) + fmt.Sprintf("%p", ch)
	r.buffer.waitQueue.Store(k, ch)
	r.closeCh()
	r.waitCh = ch
	if r.deadline.IsZero() {
		<-r.waitCh
	} else {
		err := fmt.Errorf("deadline exceeded")
		if time.Now().After(r.deadline) {
			return err
		}
		select {
		case <-r.waitCh:
		case <-time.After(r.deadline.Sub(time.Now())):
			return err
		}
	}
	return nil
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

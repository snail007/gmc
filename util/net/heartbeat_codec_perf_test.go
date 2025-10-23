// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/rand"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatCodec_Performance(t *testing.T) {
	t.Parallel()
	l, p, err := RandomListen("")
	assert.NoError(t, err)
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Server side
	go func() {
		defer wg.Done()
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		codec := NewHeartbeatCodec()
		conn := NewConn(c).AddCodec(codec)
		// Server just reads and discards data as fast as it can
		_, _ = io.Copy(io.Discard, conn)
	}()

	// Client side
	c, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)

	codec := NewHeartbeatCodec().SetInterval(time.Second)
	conn := NewConn(c).AddCodec(codec)

	dataToSend := make([]byte, 1024*4) // 4KB chunks
	_, err = rand.Read(dataToSend)
	assert.NoError(t, err)

	duration := 5 * time.Second
	startTime := time.Now()
	totalBytesSent := 0

	for time.Since(startTime) < duration {
		n, err := conn.Write(dataToSend)
		if err != nil {
			// We might get a broken pipe if the server closes connection first,
			// which is fine at the end of the test.
			break
		}
		totalBytesSent += n
	}

	elapsed := time.Since(startTime)
	rateMBs := float64(totalBytesSent) / elapsed.Seconds() / (1024 * 1024)

	t.Logf("Test finished.")
	t.Logf("Total data sent: %.2f MB", float64(totalBytesSent)/(1024*1024))
	t.Logf("Test duration: %s", elapsed.Round(time.Millisecond))
	t.Logf("Transfer rate: %.2f MB/s", rateMBs)

	// Explicitly close the client connection to signal EOF to the server.
	// This allows the server's io.Copy to complete and wg.Done() to be called.
	conn.Close()

	// Wait for server to finish
	wg.Wait()
}


// Copyright 2025 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bytes"
	"crypto/rand"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThroughputComparison(t *testing.T) {
	// Shared secret for encryption
	secret := "a very very secret key"

	// --- Test AES + Heartbeat ---
	t.Run("AES+Heartbeat", func(t *testing.T) {
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
			
			// Order: Heartbeat -> AES
			conn := NewConn(c).AddCodec(NewHeartbeatCodec()).AddCodec(NewAESCodec(secret))
			
			// Discard received data to measure pure throughput
			_, _ = io.Copy(io.Discard, conn)
		}()

		// Client side
		c, err := net.Dial("tcp", "127.0.0.1:"+p)
		assert.NoError(t, err)
		
		// Order: Heartbeat -> AES
		conn := NewConn(c).AddCodec(NewHeartbeatCodec()).AddCodec(NewAESCodec(secret))

		runThroughputTest(t, conn, "AES+Heartbeat")
		
		conn.Close()
		wg.Wait()
	})

	// --- Test GCM + Heartbeat ---
	t.Run("GCM+Heartbeat", func(t *testing.T) {
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
			
			// Order: Heartbeat -> GCM
			gcmCodec, err := NewGCMCodec(secret)
			assert.NoError(t, err)
			conn := NewConn(c).AddCodec(NewHeartbeatCodec()).AddCodec(gcmCodec)

			// Discard received data
			_, _ = io.Copy(io.Discard, conn)
		}()

		// Client side
		c, err := net.Dial("tcp", "127.0.0.1:"+p)
		assert.NoError(t, err)
		
		// Order: Heartbeat -> GCM
		gcmCodec, err := NewGCMCodec(secret)
		assert.NoError(t, err)
		conn := NewConn(c).AddCodec(NewHeartbeatCodec()).AddCodec(gcmCodec)

		runThroughputTest(t, conn, "GCM+Heartbeat")
		
		conn.Close()
		wg.Wait()
	})
}

// runThroughputTest is a helper to run the client-side data transfer and report results.
func runThroughputTest(t *testing.T, conn net.Conn, name string) {
	dataChunk := make([]byte, 1024*4) // 4KB chunks
	_, err := rand.Read(dataChunk)
	assert.NoError(t, err)

	duration := 5 * time.Second
	startTime := time.Now()
	totalBytesSent := 0
	
	var clientSentData bytes.Buffer

	for time.Since(startTime) < duration {
		n, err := conn.Write(dataChunk)
		if err != nil {
			break
		}
		clientSentData.Write(dataChunk[:n])
		totalBytesSent += n
	}

	elapsed := time.Since(startTime)
	rateMBs := float64(totalBytesSent) / elapsed.Seconds() / (1024 * 1024)

	t.Logf("--- Results for %s ---", name)
	t.Logf("Total data transferred: %.2f MB", float64(totalBytesSent)/(1024*1024))
	t.Logf("Test duration: %s", elapsed.Round(time.Millisecond))
	t.Logf("Throughput: %.2f MB/s", rateMBs)
}

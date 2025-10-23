// Copyright 2020 The GMC Author. All rights reserved.
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

func TestGCMCodec_PerformanceAndCorrectness(t *testing.T) {
	t.Parallel()
	l, p, err := RandomListen("")
	assert.NoError(t, err)
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Shared secret password for GCM key derivation
	secretPassword := "a very secret password for gcm"

	// serverReceivedData will store all data the server decrypts.
	var serverReceivedData bytes.Buffer
	// clientSentData will store all the original data the client sends.
	var clientSentData bytes.Buffer

	// Server side
	go func() {
		defer wg.Done()
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()

		// Initialize GCMCodec on the server side
		codec, err := NewGCMCodec(secretPassword)
		assert.NoError(t, err)
		conn := NewConn(c).AddCodec(codec)

		// Server reads and stores all decrypted data
		_, err = io.Copy(&serverReceivedData, conn)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF && !isUseOfClosedNetError(err) {
			t.Errorf("server copy error: %v", err)
		}
	}()

	// Client side
	c, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)

	// Initialize GCMCodec on the client side
	codec, err := NewGCMCodec(secretPassword)
	assert.NoError(t, err)
	conn := NewConn(c).AddCodec(codec)

	dataChunk := make([]byte, 1024*4) // 4KB chunks
	_, err = rand.Read(dataChunk)
	assert.NoError(t, err)

	duration := 5 * time.Second
	startTime := time.Now()
	totalBytesSent := 0

	for time.Since(startTime) < duration {
		n, err := conn.Write(dataChunk)
		if err != nil {
			break
		}
		// Store the original (unencrypted) data for later comparison
		clientSentData.Write(dataChunk)
		totalBytesSent += n
	}

	elapsed := time.Since(startTime)
	rateMBs := float64(totalBytesSent) / elapsed.Seconds() / (1024 * 1024)

	// Close the client connection to signal EOF to the server.
	conn.Close()

	// Wait for server to finish processing all data.
	wg.Wait()

	// --- Verification ---
	t.Log("Verifying data correctness...")
	assert.Equal(t, clientSentData.Len(), serverReceivedData.Len(), "Mismatch in sent and received data length")
	assert.Equal(t, clientSentData.Bytes(), serverReceivedData.Bytes(), "Mismatch in sent and received data content")
	t.Log("Data verification successful.")

	// --- Performance Results ---
	t.Logf("Total data transferred: %.2f MB", float64(totalBytesSent)/(1024*1024))
	t.Logf("Test duration: %s", elapsed.Round(time.Millisecond))
	t.Logf("Transfer rate (with AES-GCM encryption): %.2f MB/s", rateMBs)
}

// Copyright 2025 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestProtocolListener_Solution demonstrates the robust solution to the hanging issue
// by using SetFirstReadTimeout to prevent blocking on non-responsive clients.
func TestProtocolListener_Solution(t *testing.T) {
	// 1. Setup the root listener
	l, err := NewListenerAddr(":0")
	assert.NoError(t, err)
	addr := l.Addr().String()
	defer l.Close()

	// 2. THE SOLUTION: Set a short timeout for the first read.
	// This prevents the listener's accept loop from blocking on clients that connect
	// but don't send any data, which is the root cause of the hang.
	l.SetFirstReadTimeout(200 * time.Millisecond)
	l.OnFirstReadTimeout(func(ctx Context, c net.Conn, err error) {
		// This handler is called when a client is dropped due to the timeout.
		t.Logf("A client from %s timed out on first read and was dropped.", c.RemoteAddr())
	})

	// 3. Register a protocol listener for HTTP.
	// The Checker will no longer block because SetFirstReadTimeout is active.
	httpListener := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "http",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			// This PeekMax will now honor the read deadline set by SetFirstReadTimeout.
			h, err := conn.PeekMax(7)
			if err != nil {
				return false // Timeout error will be caught here
			}
			return isHTTP(h)
		},
	})

	// 4. Start an HTTP server on the httpListener
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("world"))
		})
		server := http.Server{Handler: mux}
		server.Serve(httpListener)
	}()

	// 5. Start a non-blocking loop for any other connections.
	// Although SetFirstReadTimeout handles the blocking, it's still good practice
	// to drain and close any connections that don't match a protocol.
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return // Listener closed
			}
			t.Log("Main loop accepted a non-HTTP connection and closed it.")
			conn.Close()
		}
	}()

	// 6. Connect with a simple TCP client that sends no data.
	// This client will be dropped by the FirstReadTimeout.
	t.Log("Connecting with a raw TCP client (will be dropped by timeout)...")
	tcpConn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer tcpConn.Close()

	// Wait for the timeout to occur and the client to be dropped.
	time.Sleep(300 * time.Millisecond)

	// 7. Now, make an HTTP request. It should SUCCEED because the listener is not blocked.
	t.Log("Attempting HTTP GET request...")
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s/hello", addr))

	// 8. Assert the successful outcome
	assert.NoError(t, err, "HTTP request should have succeeded, but it failed.")
	if assert.NoError(t, err) {
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		assert.NoError(t, readErr)
		assert.Equal(t, "world", string(body))
		t.Logf("HTTP request SUCCEEDED with status: %s, body: %s", resp.Status, string(body))
	}
}

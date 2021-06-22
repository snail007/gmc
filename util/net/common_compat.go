// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

// +build !go1.12,!go1.13

package gnet

import "net/http"

func bindHTTPTransport(tr *http.Transport, cc Codec, timeout time.Duration) {
	tr.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if timeout == 0 {
			timeout = time.Second * 30
		}
		c, err := Dial(addr, timeout)
		if err != nil {
			return nil, err
		}
		c.AddCodec(cc)
		return c, nil
	}
}

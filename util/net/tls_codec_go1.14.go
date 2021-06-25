// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc
// This go file used go >= go1.14.

// +build go1.14

package gnet

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
)

func (s *TLSClientCodec) getSuggestedCa(cri *tls.CertificateRequestInfo) *tls.Certificate {
	for _, chain := range s.config.Certificates {
		if err := cri.SupportsCertificate(&chain); err != nil {
			continue
		}
		return &chain
	}
	return nil
}

func (s *TLSClientCodec) setHTTPClientDialTLS(tr *http.Transport, newTLSDial func(ctx context.Context, network, addr string) (net.Conn, error)) *TLSClientCodec {
	tr.DialTLSContext = newTLSDial
	return s
}
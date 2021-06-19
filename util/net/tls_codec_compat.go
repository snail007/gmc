// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

// +build !go1.13

package gnet

func (s *TLSClientCodec) getSuggestedCa(cri *tls.CertificateRequestInfo) *tls.Certificate {
	for _, chain := range s.config.Certificates {
		if err := cri.SupportsCertificate(&chain); err != nil {
			continue
		}
		return &chain
	}
	return nil
}
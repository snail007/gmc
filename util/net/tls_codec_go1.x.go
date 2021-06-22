// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

// +build go1.13 go1.12

package gnet

import "crypto/tls"

func (s *TLSClientCodec) getSuggestedCa(*tls.CertificateRequestInfo) *tls.Certificate {
	return nil
}

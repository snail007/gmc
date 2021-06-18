// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
)

type tlsCodec struct {
	net.Conn
	config        *tls.Config
	loadSystemCas bool
	rootCas       [][]byte
}

func (s *tlsCodec) LoadSystemCas() {
	s.loadSystemCas = true
	return
}

func (s *tlsCodec) Config() *tls.Config {
	return s.config
}

func (s *tlsCodec) initRootCas(cas *x509.CertPool) (err error) {
	var c *x509.CertPool
	if s.loadSystemCas {
		c, err = x509.SystemCertPool()
		if err != nil {
			return
		}
	} else {
		c = x509.NewCertPool()
	}
	for _, v := range s.rootCas {
		if !c.AppendCertsFromPEM(v) {
			err = fmt.Errorf("parse client pem root catificate fail, %v", v)
			return
		}
	}
	*cas = *c
	return
}

func (s *tlsCodec) AddCertKeyString(certString, keyString string) (err error) {
	return s.AddCertKey([]byte(certString), []byte(keyString))
}

func (s *tlsCodec) AddCertKey(certBytes, keyBytes []byte) (err error) {
	c, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	s.addCertificate(c)
	return
}

func (s *tlsCodec) addRootCaBytes(c []byte) {
	s.rootCas = append(s.rootCas, c)
	return
}

func (s *tlsCodec) addCertificate(c tls.Certificate) {
	s.config.Certificates = append(s.config.Certificates, c)
}

type TLSClientCodec struct {
	tlsCodec
	skipVerify       bool
	pinCert          *x509.Certificate
	verifyServerName string
}

func (s *TLSClientCodec) SetVerifyServerName(verifyServerName string) {
	s.verifyServerName = verifyServerName
}

func (s *TLSClientCodec) PinVerifyCert(serverCertPEMBytes []byte) (err error) {
	block, _ := pem.Decode(serverCertPEMBytes)
	if block == nil {
		err = fmt.Errorf("failed to parse certificate PEM")
		return
	}
	s.pinCert, err = x509.ParseCertificate(block.Bytes)
	return
}

func NewTLSClientCodec(c *tls.Config) *TLSClientCodec {
	if c.RootCAs == nil {
		c.RootCAs = x509.NewCertPool()
	}
	c.InsecureSkipVerify = true
	s := &TLSClientCodec{
		tlsCodec: tlsCodec{
			config: c,
		},
	}
	return s
}

func (s *TLSClientCodec) Initialize(ctx Context, next NextCodec) (c net.Conn, err error) {
	s.initRootCas(s.config.RootCAs)
	if !s.skipVerify {
		s.config.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			verifyOkay := false
			if s.pinCert != nil {
				for _, rawCert := range rawCerts {
					cert, _ := x509.ParseCertificate(rawCert)
					if s.pinCert.Equal(cert) {
						verifyOkay = true
					}
				}
			} else {
				for _, rawCert := range rawCerts {
					cert, _ := x509.ParseCertificate(rawCert)
					_, e := cert.Verify(x509.VerifyOptions{
						DNSName: s.verifyServerName,
						Roots:   s.config.RootCAs,
					})
					if e == nil {
						verifyOkay = true
						break
					}
				}
			}
			if !verifyOkay {
				return fmt.Errorf("unknown server certificate received")
			}
			return nil
		}
		if s.verifyServerName == "" || s.pinCert != nil {

		}
	}
	s.Conn = tls.Client(ctx.Conn(), s.config)
	return next.Call(ctx.SetConn(s))
}

func (s *TLSClientCodec) AddVerifyCa(caPEMBytes []byte) {
	s.addRootCaBytes(caPEMBytes)
	return
}

func (s *TLSClientCodec) SkipVerify(b bool) {
	s.skipVerify = b
}

type TLSServerCodec struct {
	tlsCodec
}

func NewTLSServerCodec(c *tls.Config) *TLSServerCodec {
	if c.ClientCAs == nil {
		c.ClientCAs = x509.NewCertPool()
	}
	s := &TLSServerCodec{
		tlsCodec: tlsCodec{
			config: c,
		},
	}
	return s
}

func (s *TLSServerCodec) Initialize(ctx Context, next NextCodec) (c net.Conn, err error) {
	s.initRootCas(s.config.ClientCAs)
	s.Conn = tls.Server(ctx.Conn(), s.config)
	return next.Call(ctx.SetConn(s))
}

func (s *TLSServerCodec) RequireClientAuth() {
	s.config.ClientAuth = tls.RequireAndVerifyClientCert
}

func (s *TLSServerCodec) AddClientCa(caPEMBytes []byte) {
	s.addRootCaBytes(caPEMBytes)
	return
}

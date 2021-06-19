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
	"strings"
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

func (s *tlsCodec) AddCertificate(certPEMBytes, keyPEMBytes []byte) (err error) {
	c, err := tls.X509KeyPair(certPEMBytes, keyPEMBytes)
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
	skipVerify           bool
	pinCert              *x509.Certificate
	serverName           string
	skipVerifyCommonName bool
}

func (s *TLSClientCodec) SkipVerifyCommonName(skipVerifyCommonName bool) {
	s.skipVerifyCommonName = skipVerifyCommonName
}

func (s *TLSClientCodec) SetServerName(serverName string) {
	s.serverName = serverName
}

func (s *TLSClientCodec) PinServerCert(serverCertPEMBytes []byte) (err error) {
	block, _ := pem.Decode(serverCertPEMBytes)
	if block == nil {
		err = fmt.Errorf("failed to parse certificate PEM")
		return
	}
	s.pinCert, err = x509.ParseCertificate(block.Bytes)
	return
}

func (s *TLSClientCodec) Initialize(ctx Context, next NextCodec) (c net.Conn, err error) {
	if s.config.RootCAs == nil {
		s.config.RootCAs = x509.NewCertPool()
	}
	err = s.initRootCas(s.config.RootCAs)
	if err != nil {
		return nil, err
	}
	s.config.InsecureSkipVerify = true
	s.config.ServerName = s.serverName
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
					opts := x509.VerifyOptions{
						// verify server's certificate chains must be signed by s.config.RootCAs.
						Roots: s.config.RootCAs,
					}
					if !s.skipVerifyCommonName {
						// verify server's certificate chains Common Name must be contains s.ServerName.
						opts.DNSName = s.serverName
						if cert.DNSNames == nil {
							// cert.Verify will only check servername from cert.DNSNames, without cert.Subject.CommonName,
							// so we append the serverName into cert.DNSNames when cert.DNSNames is nil, ensure Verify success.
							cert.DNSNames = append(cert.DNSNames, s.serverName)
						}
					}
					_, e := cert.Verify(opts)
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
	}
	s.config.GetClientCertificate = func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		// only one certificate, just using it.
		if len(s.config.Certificates) == 1 {
			return &s.config.Certificates[0], nil
		}

		// find by request server name.
		for _, chain := range s.config.Certificates {
			cert, _ := x509.ParseCertificate(chain.Certificate[0])
			// exactly equal
			if cert.Subject.CommonName == s.serverName {
				return &chain, nil
			}
			// exactly equal in DNSNames
			for _, name := range cert.DNSNames {
				if s.serverName == name {
					return &chain, nil
				}
			}
			// wildcard name equal in DNSNames
			labels := strings.Split(s.serverName, ".")
			labels[0] = "*"
			wildcardName := strings.Join(labels, ".")
			for _, name := range cert.DNSNames {
				if wildcardName == name {
					return &chain, nil
				}
			}
		}

		// find by server response supported root cas.
		for _, chain := range s.config.Certificates {
			if err := cri.SupportsCertificate(&chain); err != nil {
				continue
			}
			return &chain, nil
		}

		// No acceptable certificate found. Don't send a certificate.
		return new(tls.Certificate), nil
	}
	s.Conn = tls.Client(ctx.Conn(), s.config)
	return next.Call(ctx.SetConn(s))
}

func (s *TLSClientCodec) AddServerCa(caPEMBytes []byte) {
	s.addRootCaBytes(caPEMBytes)
	return
}

func (s *TLSClientCodec) SkipVerify(b bool) {
	s.skipVerify = b
}

func NewTLSClientCodec() *TLSClientCodec {
	s := &TLSClientCodec{
		tlsCodec: tlsCodec{
			config: &tls.Config{},
		},
	}
	return s
}

type TLSServerCodec struct {
	tlsCodec
}

func (s *TLSServerCodec) Initialize(ctx Context, next NextCodec) (c net.Conn, err error) {
	if s.config.ClientCAs == nil {
		s.config.ClientCAs = x509.NewCertPool()
	}
	err = s.initRootCas(s.config.ClientCAs)
	if err != nil {
		return nil, err
	}
	// To make sure choose the server certificate depends on the request server name.
	s.config.BuildNameToCertificate()
	s.Conn = tls.Server(ctx.Conn(), s.config)
	return next.Call(ctx.SetConn(s))
}

func (s *TLSServerCodec) RequireClientAuth(b bool) {
	if b {
		s.config.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		s.config.ClientAuth = tls.NoClientCert
	}
}

func (s *TLSServerCodec) AddClientCa(caPEMBytes []byte) {
	s.addRootCaBytes(caPEMBytes)
	return
}

func NewTLSServerCodec() *TLSServerCodec {
	s := &TLSServerCodec{
		tlsCodec: tlsCodec{
			config: &tls.Config{},
		},
	}
	return s
}

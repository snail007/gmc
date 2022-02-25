// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"hash"
	"io"
	"net"

	"golang.org/x/crypto/pbkdf2"
)

const (
	aesDefaultIterations = 4096
	aesDefaultKeySize    = 32 //256bits
)

var (
	aesDefaultHashFunc = sha256.New
	aesDefaultSalt     = []byte(`
V6Nt!d|@bo+N$L9+<d$|(;QUHj.BQ?RXzYSO]ifkXp/G!kFmWyXyEg6e26T}
VA-wxBptOTM^2,id,6acYKY_ecP5%)wnAo<:>SO+(x"R";\'4&fTAVu92GhW
Snsgymt!3gbP2pe=J//}1a?lp9ej=&TB!C_V(cT2?z8wyoL_-1hwk=]3fd[]
`)
	aesNewCipher = aes.NewCipher
)

type AESOptions struct {
	// Password required.
	Password string
	// Salt optionally. A good random salt at least 8 bytes is recommended.
	Salt []byte
	// Type optionally. It must be: aes-128 or aes-192 or aes-256, if empty default is aes-256.
	Type string
	// HashFunc optionally, default is sha256.New
	HashFunc func() hash.Hash
}

type AESCodec struct {
	net.Conn
	password string
	key      []byte
	r        io.Reader
	w        io.Writer
}

func (s *AESCodec) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}

func (s *AESCodec) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

func (s *AESCodec) Close() error {
	if s.Conn != nil {
		return s.Conn.Close()
	}
	return nil
}

func (s *AESCodec) Initialize(ctx Context, next NextCodec) (conn net.Conn, err error) {
	block, err := aesNewCipher(s.key)
	if err != nil {
		return
	}
	s.Conn = ctx.Conn()
	riv := sha256.New().Sum(s.key)
	rStream := cipher.NewCFBDecrypter(block, riv[:aes.BlockSize])
	s.r = &cipher.StreamReader{S: rStream, R: s.Conn}
	wiv := sha256.New().Sum(riv)
	wStream := cipher.NewCFBEncrypter(block, wiv[:aes.BlockSize])
	s.w = &cipher.StreamWriter{S: wStream, W: s.Conn}
	return next.Call(ctx.SetConn(s))
}

func NewAESCodec(password string) *AESCodec {
	c := new(AESOptions)
	c.Password = password
	return NewAESCodecFromOptions(c)
}

func NewAESCodecFromOptions(c *AESOptions) *AESCodec {
	if len(c.Salt) == 0 {
		c.Salt = aesDefaultSalt
	}

	iterations := aesDefaultIterations
	keySize := aesDefaultKeySize
	if c.Type != "" {
		switch c.Type {
		case "aes-128":
			iterations = 1024
			keySize = 16
		case "aes-192":
			iterations = 2048
			keySize = 24
		case "aes-256":
			iterations = 4096
			keySize = 32
		}
	}
	if c.HashFunc == nil {
		c.HashFunc = aesDefaultHashFunc
	}
	key := pbkdf2.Key([]byte(c.Password), c.Salt, iterations, keySize, c.HashFunc)
	return &AESCodec{key: key}
}

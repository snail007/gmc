// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

//go:build !linux
// +build !linux

package redir

import (
	"net"
)

func RealServerAddress(conn net.Conn) (string, error) {
	return "", nil
}

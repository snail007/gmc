// +build !linux

package redir

import (
	"net"
)

func RealServerAddress(conn net.Conn) (string, error) {
	return "", nil
}

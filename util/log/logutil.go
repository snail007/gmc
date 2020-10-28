package logutil

import (
	gmclog "github.com/snail007/gmc/base/log"
	gmccore "github.com/snail007/gmc/core"
)

func New(prefix string) gmccore.Logger {
	l := gmclog.NewGMCLog()
	if prefix != "" {
		l.With(prefix)
	}
	return l
}


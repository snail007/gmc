package logutil

import (
	"log"
	"os"
)

func New(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds)
}
func NewStd(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}

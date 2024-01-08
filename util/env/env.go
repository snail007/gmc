package genv

import (
	gvalue "github.com/snail007/gmc/util/value"
	"os"
)

// Accessor access to os environment variable with auto prefix support
// and the value is gvalue.AnyValue can be used as many data type
type Accessor struct {
	prefix string
}

func NewAccessor(prefix string) *Accessor {
	return &Accessor{prefix: prefix}
}

func (s *Accessor) addPrefix(str string) string {
	return s.prefix + str
}

func (s *Accessor) Get(key string) *gvalue.AnyValue {
	v := os.Getenv(s.addPrefix(key))
	return gvalue.NewAny(v)
}

func (s *Accessor) Lookup(key string) (*gvalue.AnyValue, bool) {
	v, found := os.LookupEnv(s.addPrefix(key))
	if !found {
		return nil, false
	}
	return gvalue.NewAny(v), true
}

func (s *Accessor) Set(key, value string) *Accessor {
	os.Setenv(s.addPrefix(key), value)
	return s
}

func (s *Accessor) Unset(key string) *Accessor {
	os.Unsetenv(s.addPrefix(key))
	return s
}

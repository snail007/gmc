package gcond

import (
	"github.com/snail007/gmc/util/value"
)

func Cond(check bool, ok interface{}, fail interface{}) *gvalue.Value {
	if check {
		return gvalue.New(ok)
	}
	return gvalue.New(fail)
}

func CondFn(check bool, ok func() interface{}, fail func() interface{}) *gvalue.Value {
	if check {
		return gvalue.New(ok())
	}
	return gvalue.New(fail())
}

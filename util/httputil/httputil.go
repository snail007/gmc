package gmchttputil

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func Die(w io.Writer, data ...interface{}) {
	Write(w, data...)
	panic("__DIE__")
}

func Stop(w io.Writer, data ...interface{}) {
	Write(w, data...)
	panic("__STOP__")
}

// StopE will exit controller method if error is not nil.
// First argument is an error.
// Secondary argument is fail function, it be called if error is not nil.
// Third argument is success function, it be called if error is nil.
func StopE(err interface{},fn  ...func()) {
	var failFn, okayFn func()
	if len(fn)>=1{
		failFn=fn[0]
	}
	if len(fn)>=2{
		okayFn=fn[1]
	}
	if err==nil{
		if okayFn!=nil{
			okayFn()
		}
		return
	}
	if failFn!=nil{
		failFn()
	}
	panic("__STOP__")
}
func Write(w io.Writer, data ...interface{}) (n int, err error) {
	for _, d := range data {
		if d == nil {
			continue
		}
		switch v := d.(type) {
		case []byte:
			n, err = w.Write(v)
		case string:
			n, err = w.Write([]byte(v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			n, err = w.Write([]byte(fmt.Sprintf("%d", v)))
		case bool:
			str := "true"
			if !v {
				str = "false"
			}
			n, err = w.Write([]byte(str))
		case float32:
			n, err = w.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		case float64:
			n, err = w.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
		case error:
			n, err = w.Write([]byte(v.Error()))
		default:
			t := reflect.TypeOf(v)
			//map, slice
			//jsonType := []string{"[", "map["}
			found := false
			//vTypeStr := t.String()
			//for _, typ := range jsonType {
			if t.Kind()==reflect.Slice || t.Kind()==reflect.Map {
				found = true
				var b []byte
				b, err = json.Marshal(v)
				if err == nil {
					n, err = w.Write(b)
				}
				//break
			}
			//}
			if !found {
				fmt.Println(found)
				n, err = w.Write([]byte(fmt.Sprintf("unsupported type to write: %s", t.String())))
			}
		}
		if err != nil {
			return
		}
	}
	return
}

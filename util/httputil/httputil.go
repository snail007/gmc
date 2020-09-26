package gmchttputil

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
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
func Write(w io.Writer, data ...interface{}) (err error) {
	for _, d := range data {
		if d == nil {
			continue
		}
		switch v := d.(type) {
		case []byte:
			w.Write(v)
		case string:
			w.Write([]byte(v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			w.Write([]byte(fmt.Sprintf("%d", v)))
		case bool:
			str := "true"
			if !v {
				str = "false"
			}
			w.Write([]byte(str))
		case float32, float64:
			w.Write([]byte(fmt.Sprintf("%f", v)))
		case error:
			w.Write([]byte(v.Error()))
		default:
			t := reflect.TypeOf(v)
			//map, slice
			jsonType := []string{"[", "map["}
			found := false
			vTypeStr := t.String()
			for _, typ := range jsonType {
				if strings.HasPrefix(vTypeStr, typ) {
					found = true
					var b []byte
					b, err = json.Marshal(v)
					if err == nil {
						w.Write(b)
					}
					break
				}
			}
			if !found {
				w.Write([]byte(fmt.Sprintf("unsupported type to write: %s", t.String())))
			}
		}
	}
	return
}

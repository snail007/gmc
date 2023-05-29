package greflect

import (
	"reflect"
)

func DeepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()

	switch srcType.Kind() {
	case reflect.Ptr:
		if srcValue.IsNil() {
			return nil
		}
		elem := srcValue.Elem()
		dest := reflect.New(elem.Type())
		dest.Elem().Set(reflect.ValueOf(DeepCopy(elem.Interface())))
		return dest.Interface()
	default:
		dest := reflect.New(srcType).Elem()
		deepCopyValue(srcValue, dest)
		return dest.Interface()
	}
}

func deepCopyValue(src, dest reflect.Value) {
	srcType := src.Type()

	switch srcType.Kind() {
	case reflect.Slice:
		dest.Set(reflect.MakeSlice(srcType, src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			elem := src.Index(i)
			copyElem := reflect.New(elem.Type()).Elem()
			deepCopyValue(elem, copyElem)
			dest.Index(i).Set(copyElem)
		}
	case reflect.Map:
		dest.Set(reflect.MakeMap(srcType))
		for _, key := range src.MapKeys() {
			elem := src.MapIndex(key)
			copyElem := reflect.New(elem.Type()).Elem()
			deepCopyValue(elem, copyElem)
			dest.SetMapIndex(key, copyElem)
		}
	case reflect.Ptr:
		elem := src.Elem()
		dest.Set(reflect.New(elem.Type()))
		deepCopyValue(elem, dest.Elem())
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			field := src.Field(i)
			fieldType := srcType.Field(i)
			if fieldType.PkgPath == "" {
				copyField := reflect.New(field.Type()).Elem()
				deepCopyValue(field, copyField)
				dest.Field(i).Set(copyField)
			}
		}
	default:
		dest.Set(src)
	}
}

package util

import (
	"errors"
	"reflect"
)

var ErrCopyType error = errors.New("copy type error")

func CopyStruct(dst, src interface{}) (err error) {
	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
		if srcType.Kind() != reflect.Struct {
			err = ErrCopyType
			return
		}
		srcValue = srcValue.Elem()
	}

	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		err = ErrCopyType
		return
	}

	dstType = dstType.Elem()
	dstValue := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstType.NumField(); i++ {
		dv := dstValue.Field(i)
		if !dv.IsValid() || !dv.CanSet() {
			continue
		}

		sv := srcValue.FieldByName(dstType.Field(i).Name)
		if !sv.IsValid() {
			continue
		}
		if sv.Type() != dv.Type() {
			continue
		}
		dv.Set(sv)
	}
	return
}

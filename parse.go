package clop

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotPointerType = errors.New("Not pointer type")
var ErrUnsupportedType = errors.New("Unsupported type")

func parseStruct(x interface{}) error {
	v := reflect.ValueOf(x)

	if x == nil || v.IsNil() {
		return ErrUnsupportedType
	}

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("%s:got(%T)", ErrNotPointerType, v.Type())
	}

	for {
		v = v.Elem()
		if v.Kind() != reflect.Ptr {
			break
		}
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)
		fmt.Printf("%s\n", sf.Tag)
		fmt.Printf("index(%d)(1.%s)-->(2.%s)\n", i, sf.Tag.Get("clop"), sf.Tag.Get("usage"))
	}

	return nil
}

package clop

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type convert struct {
	bitSize int
	cb      func(val string, bitSize int, field reflect.Value) error
}

var convertFunc map[reflect.Kind]convert

func init() {
	convertFunc = map[reflect.Kind]convert{
		reflect.Uint:    {bitSize: 0, cb: setUintField},
		reflect.Uint8:   {bitSize: 8, cb: setUintField},
		reflect.Uint16:  {bitSize: 16, cb: setUintField},
		reflect.Uint32:  {bitSize: 32, cb: setUintField},
		reflect.Uint64:  {bitSize: 64, cb: setUintField},
		reflect.Int:     {bitSize: 0, cb: setIntField},
		reflect.Int8:    {bitSize: 8, cb: setIntField},
		reflect.Int16:   {bitSize: 16, cb: setIntField},
		reflect.Int32:   {bitSize: 32, cb: setIntField},
		reflect.Int64:   {bitSize: 64, cb: setIntDurationField},
		reflect.Bool:    {bitSize: 0, cb: setBoolField},
		reflect.Float32: {bitSize: 32, cb: setFloatField},
		reflect.Float64: {bitSize: 64, cb: setFloatField},
		reflect.Struct:  {bitSize: 0, cb: setStructField},
		reflect.Slice:   {bitSize: 0, cb: setSlice},
		reflect.Map:     {bitSize: 0, cb: setMapField},
	}
}

func setIntDurationField(val string, bitSize int, value reflect.Value) error {
	switch value.Interface().(type) {
	case time.Duration:
		return setTimeDuration(val, bitSize, value)
	}

	return setIntField(val, bitSize, value)
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}

	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setStructField(val string, bitSize int, value reflect.Value) error {
	switch value.Interface().(type) {
	case time.Time:
		//return setTimeField(val, bitSize, value)
	}
	return json.Unmarshal([]byte(val), value.Addr().Interface())
}

func setSlice(val string, bitSize int, value reflect.Value) error {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	first := false
	if value.Len() == 0 {
		first = true
		// 初始化一个不为空的slice
		value.Set(reflect.MakeSlice(value.Type(), 1, 1))
	}

	v := reflect.New(value.Index(0).Type())

	v = v.Elem()
	if err := setBase(val, v); err != nil {
		return err
	}

	if first {
		value.SetLen(0)
	}
	v2 := reflect.Append(value, v)
	value.Set(v2)

	return nil
}

func setMapField(val string, bitSize int, value reflect.Value) error {
	return json.Unmarshal([]byte(val), value.Addr().Interface())
}

func setTimeDuration(val string, bitSize int, value reflect.Value) error {
	if val == "" {
		val = "0"
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func setBase(val string, value reflect.Value) error {
	if value.Kind() == reflect.String {
		value.SetString(val)
		return nil
	}

	fn, ok := convertFunc[value.Kind()]
	if ok {
		return fn.cb(val, fn.bitSize, value)
	}

	return fmt.Errorf("type (%T) unsupported type", value)
}

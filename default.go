package clop

import (
	"encoding/json"
	"reflect"
)

func isDefvalJSON(def string) bool {
	bs := StringToBytes(def)
	return json.Valid(bs)
}

func setDefaultValue(def string, v reflect.Value) error {
	if isDefvalJSON(def) {
		return json.Unmarshal([]byte(def), v.Addr().Interface())
	}

	return setBase(def, v)
}

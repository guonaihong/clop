package clop

import (
	"encoding/json"
	"reflect"
)

func isDefvalJSON(def []byte) bool {
	if def[0] == '[' && len(def) > 2 && def[len(def)-1] == ']' ||
		def[0] == '{' && len(def) > 2 && def[len(def)-1] == '}' {

		return json.Valid(def)
	}
	return false
}

func setDefaultValue(def string, v reflect.Value) error {
	def2 := StringToBytes(def)
	if isDefvalJSON(def2) {
		return json.Unmarshal(def2, v.Addr().Interface())
	}

	return setBase(def, v)
}

package clop

import "strings"

// gnu风格名字是，下划线(蛇形)或者(驼峰)都换成中横线风格
// LongOpt -> short-opt
// long_opt -> long-opt

func wordStart(b byte) bool {

	return b >= 'A' && b <= 'Z' || b == '_'
}

func gnuOptionName(opt string) (string, error) {

	var name strings.Builder

	for i, b := range []byte(opt) {

		if wordStart(b) {
			if i != 0 {
				name.WriteByte('-')
			}

			if b != '_' {
				b = b - 'A' + 'a'
				name.WriteByte(b)
			}

			continue
		}

		name.WriteByte(b)

	}

	return name.String(), nil
}

func envOptionName(opt string) (string, error) {

	var name strings.Builder

	for i, b := range []byte(opt) {

		if wordStart(b) {
			if i != 0 {
				name.WriteByte('_')
			}

			if b == '_' {
				continue
			}

			name.WriteByte(b)
			continue
		}

		if b >= 'a' && b <= 'z' {
			name.WriteByte(b - 'a' + 'A')
		} else {
			name.WriteByte('_')
		}

	}

	return name.String(), nil
}

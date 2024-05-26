package clop

import "strings"

// gnu风格名字是，下划线(蛇形)或者(驼峰)都换成中横线风格
// LongOpt -> short-opt
// long_opt -> long-opt

// 专有名字不转换
// 专有名字词
var specialNames = map[string]bool{
	"JSON": true,
	"XML":  true,
	"YAML": true,
	"URL":  true,
	"URI":  true,
}

func wordStart(b byte) bool {

	return b >= 'A' && b <= 'Z' || b == '_'
}

// gnuOptionName 转换为gnu风格的名字
func gnuOptionName(opt string) (string, error) {

	var name strings.Builder

	if specialNames[opt] {
		return opt, nil
	}

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

// 环境变量名字是大写，下划线(蛇形)风格
func envOptionName(opt string) (string, error) {

	var name strings.Builder

	if specialNames[opt] {
		return opt, nil
	}
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

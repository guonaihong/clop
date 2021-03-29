package clop

import "strings"

// gnu风格名字是，下划线(蛇形)或者(驼峰)都换成中横线风格
// LongOpt -> short-opt
// long_opt -> long-opt

func gnuOptionName(opt string) (string, error) {

	var name strings.Builder

	for i, b := range []byte(opt) {

		if b >= 'A' && b <= 'Z' || b == '_' {
			if i != 0 {
				name.WriteByte('-')
			}

			if b != '_' {
				b = b - 'A' + 'a'
				err := name.WriteByte(b)
				if err != nil {
					return "", err
				}
			}

			continue
		}

		name.WriteByte(b)

	}

	return name.String(), nil
}

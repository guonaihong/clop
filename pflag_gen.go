package clop

import (
	"bytes"
	"fmt"
)

// 根据解析的函数名和参数, 生成结构体
func genStructBytes(p *ParseFlag) {

	var code bytes.Buffer
	for k, funcAndArgs := range p.funcAndArgs {
		v := funcAndArgs
		if !v.haveParse {
			continue
		}

		code.WriteString(fmt.Sprintf("type %s struct{", k))

		for _, arg := range v.args {
			// 选项名是比较重要的, 没有就不生成
			if len(arg.optName) == 0 {
				continue
			}
			// 写入字段名和类型名
			code.WriteString(fmt.Sprintf("%s %s", arg.varName, arg.typeName))

			var clopTag bytes.Buffer

			// 写入选项名
			if len(arg.optName) > 0 {
				clopTag.WriteString("`clop:\"")
				numMinuses := "-"
				if len(arg.optName) > 1 {
					numMinuses = "--"
				}
				code.WriteString(fmt.Sprintf("%s%s\"", numMinuses, arg.optName))
			}

			// 写入默认值
			if len(arg.defVal) > 0 {
				clopTag.WriteString(fmt.Sprintf("default:\"%s\"", arg.defVal))
			}

			// 写入帮助信息
			if len(arg.usage) > 0 {
				clopTag.WriteString(fmt.Sprintf("usage:\"%s\"", arg.usage))
			}

		}

		code.WriteString("}")

		p.allOutBuf.Write(code.Bytes())

		code.Reset()

	}

}

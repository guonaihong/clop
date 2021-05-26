package clop

import (
	"bytes"
	"fmt"
	"go/format"
)

// 根据解析的函数名和参数, 生成结构体
func genStructBytes(p *ParseFlag) error {

	var code bytes.Buffer
	for k, funcAndArgs := range p.funcAndArgs {
		v := funcAndArgs
		if !v.haveParseFunc {
			continue
		}

		code.WriteString(fmt.Sprintf("type %sAutoGen struct{", k))

		for _, arg := range v.args {
			// 选项名是比较重要的, 没有就不生成
			if len(arg.optName) == 0 || len(arg.varName) == 0 {
				continue
			}
			// 写入字段名和类型名
			varName := arg.varName
			if varName[0] >= 'a' && varName[0] <= 'z' {
				varName = string(varName[0]-'a'+'A') + varName[1:]
			}

			code.WriteString(fmt.Sprintf("%s %s", varName, arg.typeName))

			var clopTag bytes.Buffer

			// 写入选项名
			if len(arg.optName) > 0 {
				clopTag.WriteString("`clop:\"")
				numMinuses := "-"
				if len(arg.optName) > 1 {
					numMinuses = "--"
				}
				clopTag.WriteString(fmt.Sprintf("%s%s\" ", numMinuses, arg.optName))
			}

			// 写入默认值
			if len(arg.defVal) > 0 {
				clopTag.WriteString(fmt.Sprintf("default:\"%s\" ", arg.defVal))
			}

			// 写入帮助信息
			if len(arg.usage) > 0 {
				clopTag.WriteString(fmt.Sprintf("usage:\"%s\" `\n", arg.usage))
			}

			code.WriteString(clopTag.String())

		}

		code.WriteString("}")

		p.allOutBuf.Write(code.Bytes())

		code.Reset()

	}

	fmtCode, err := format.Source(p.allOutBuf.Bytes())
	if err != nil {
		return err
	}

	p.allOutBuf.Reset()
	p.allOutBuf.Write(fmtCode)

	return nil
}

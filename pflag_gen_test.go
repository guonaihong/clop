package clop

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenStructBytes(t *testing.T) {

	p := ParseFlag{

		funcAndArgs: map[string]funcAndArgs{
			"flag": funcAndArgs{
				args: []flagOpt{
					{
						varName:  "header",
						optName:  "header",
						defVal:   "",
						usage:    "test header usage",
						typeName: "string",
					},
					{
						varName:  "size",
						optName:  "size",
						defVal:   "",
						usage:    "test size usage",
						typeName: "int",
					},
				},
				haveParseFunc: true,
			},
		},
	}

	err := genStructBytes(&p)
	assert.NoError(t, err)
	io.Copy(os.Stdout, &p.allOutBuf)
}

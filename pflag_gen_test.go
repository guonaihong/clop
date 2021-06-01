package clop

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenStructBytes(t *testing.T) {

	p := ParseFlag{

		haveStruct:     true,
		haveImportPath: true,
		haveMain:       true,
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
			"flag2": funcAndArgs{
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

	all, err := genStructBytes(&p)
	assert.NoError(t, err)
	os.Stdout.Write(all)
}

package clop

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Pflag_Parser_Case struct {
	fileName string
	needLine int
}

func Test_Pflag_Parse(t *testing.T) {

	testcase := []Pflag_Parser_Case{
		{"./testdata/flag.go.tst", 3},
		{"./testdata/flag_ptr.go.tst", 5},
	}

	for _, tc := range testcase {

		p := NewParseFlag().FromFile(tc.fileName)

		all, err := p.Parse()
		assert.NoError(t, err)

		fmt.Printf("%#v\n", p.funcAndArgs)
		fmt.Printf("%s\n", string(all))
		countLine := bytes.Count(all, []byte("\n"))
		//fmt.Printf("%d\n", countLine)

		assert.Equal(t, tc.needLine, countLine)
	}

	/*
		for k, v := range p.funcAndArgs {
			fmt.Printf("k(%s), v(%v)\n", k, v)
		}
	*/
}

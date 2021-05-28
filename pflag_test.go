package clop

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pflag_Parse(t *testing.T) {
	p := NewParseFlag().FromFile("./testdata/flag.go.tst")
	all, err := p.Parse()
	assert.NoError(t, err)

	fmt.Printf("%s\n", string(all))
	/*
		for k, v := range p.funcAndArgs {
			fmt.Printf("k(%s), v(%v)\n", k, v)
		}
	*/
}

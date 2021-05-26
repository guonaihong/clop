package clop

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pflag_Parse(t *testing.T) {
	p := NewParseFlag().FromFile("./testdata/flag.go")
	all, err := p.Parse()
	assert.NoError(t, err)

	fmt.Printf("%s\n", string(all))
}

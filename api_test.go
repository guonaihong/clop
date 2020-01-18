package clop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_API_test(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"short=c;long=number-nonblank"
							usage:"number nonempty output lines, overrides"`

		ShowEnds bool `clop:"short=E;long=show-ends" 
							usage:"display $ at end of each line"`

		Number bool `clop:"short=n;long=number" 
							usage:"number all output lines"`

		SqueezeBlank bool `clop:"short=n;long=squeeze-blank"
							usage:"suppress repeated empty output lines"`

		ShowTab bool `clop:"short=s;long=show-tabs"
							usage:"display TAB characters as ^I"`

		ShowNonprinting bool `clop:"short=v;long=show-nonprinting"
							usage:"use ^ and M- notation, except for LFD and TAB" `
	}

	c := cat{}
	err := parseStruct(&c)
	assert.NoError(t, err)
}

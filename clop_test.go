package clop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_API_cat_test(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-c;--number-nonblank"
							usage:"number nonempty output lines, overrides"`

		ShowEnds bool `clop:"-E;--show-ends"
						    usage:"display $ at end of each line"`

		Number bool `clop:"-n;--number"
						    usage:"number all output lines"`

		SqueezeBlank bool `clop:"-n;--squeeze-blank"
						    usage:"suppress repeated empty output lines"`

		ShowTab bool `clop:"-s;--show-tabs"
						    usage:"display TAB characters as ^I"`

		ShowNonprinting bool `clop:"-v;--show-nonprinting"
						    usage:"use ^ and M- notation, except for LFD and TAB" `
	}

	c := cat{}

	cp := New([]string{})
	err := cp.register(&c)
	assert.NoError(t, err)
}

func Test_API_head_test(t *testing.T) {
	type head struct {
		Bytes int `clop:"-c;--bytes"
				   usage:" print the first NUM bytes of each file;
	      		     with the leading '-', print all but the last
		  		     NUM bytes of each file"`

		Lines int `clop:"-n;--lines;-\d+,regex"
				   usage:"print the first NUM lines instead of the first 10;
                     with the leading '-', print all but the last
                     NUM lines of each file"`

		Quiet bool `clop:"-q;--quiet;--silent"
				   usage:"never print headers giving file names"`

		Verbose bool `clop:"-v;--verbose"
				   usage:"always print headers giving file names"`

		ZeroTerminated byte `clop:"-z;--zero-terminated;def='\n'" 
							usage:"line delimiter is NUL, not newline"`
	}

	h := head{}
	cp := New([]string{})
	err := cp.register(&h)
	assert.NoError(t, err)
}

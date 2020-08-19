package clop

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenMaybeHelpMsg(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-b;--number-nonblank"
		                     usage:"number nonempty output lines, overrides"`

		ShowEnds bool `clop:"-E;--show-ends"
		               usage:"display $ at end of each line"`

		Number bool `clop:"-n;--number"
		             usage:"number all output lines"`

		SqueezeBlank bool `clop:"-s;--squeeze-blank"
		                   usage:"suppress repeated empty output lines"`

		ShowTab bool `clop:"-T;--show-tabs"
		              usage:"display TAB characters as ^I"`

		ShowNonprinting bool `clop:"-v;--show-nonprinting"
		                      usage:"use ^ and M- notation, except for LFD and TAB" `

		Files []string `clop:"args=files"`
	}

	var out strings.Builder
	for index, test := range []testAPI{
		// 测试--number-nonblank 写错的可能
		{
			func() cat {
				c := cat{}
				cp := New([]string{"--num-nonblank"}).SetExit(false).SetOutput(&out)
				err := cp.Bind(&c)
				assert.Error(t, err)
				assert.NotEqual(t, -1, strings.Index(out.String(), "--number-nonblank"))
				out.Reset()
				return c
			}(), cat{},
		},
		// 测试--show-ends 写错的可能
		{
			func() cat {
				c := cat{}
				cp := New([]string{"--show-end"}).SetExit(false).SetOutput(&out)
				err := cp.Bind(&c)
				assert.Error(t, err)
				assert.NotEqual(t, -1, strings.Index(out.String(), "--show-ends"))
				out.Reset()
				return c
			}(), cat{},
		},
	} {

		assert.Equal(t, test.got, test.need, fmt.Sprintf("index = %d", index))
	}
}

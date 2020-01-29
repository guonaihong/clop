package clop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAPI struct {
	got  interface{}
	need interface{}
}

func Test_API_cat_test(t *testing.T) {
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

		Args []string `clop:"args"`
	}

	for _, test := range []testAPI{
		// 测试短选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-vTsnEb"})
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},
		// 测试长选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"--show-nonprinting", "--show-tabs", "--squeeze-blank", "--number", "--show-ends", "--number-nonblank"})
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},

		// 测试长短选项混合的情况
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-v", "--show-tabs", "-s", "--number", "-E", "--number-nonblank"})
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},

		// 测试args选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-n", "r.go", "-T", "pool.c"})
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{Number: true, ShowTab: true, Args: []string{"r.go", "pool.c"}},
		},
	} {

		assert.Equal(t, test.got, test.need)
	}
}

func Test_API_grep_test(t *testing.T) {
	type grep struct {
		BeforeContext int      `clop:"-B;--before-context" usage:"print NUM lines of leading context"`
		AfterContext  int      `clop:"-A;--after-context"   usage:"print NUM lines of trailing context"`
		Args          []string `clop:"args"`
	}

	for _, test := range []testAPI{
		// 测试短选项
		{
			func() grep {
				g := grep{}
				cp := New([]string{"-B3", "--after-context", "1", "keyword", "join.txt"})
				err := cp.Bind(&g)
				assert.NoError(t, err)
				return g
			}(), grep{BeforeContext: 3, AfterContext: 1, Args: []string{"keyword", "join.txt"}},
		},
	} {
		assert.Equal(t, test.need, test.got)
	}
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

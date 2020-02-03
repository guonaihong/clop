package clop

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testAPI struct {
	got  interface{}
	need interface{}
}

// 测试bool类型
func Test_API_bool(t *testing.T) {
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

	for index, test := range []testAPI{
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

		assert.Equal(t, test.got, test.need, fmt.Sprintf("index = %d", index))
	}
}

// 测试slice
func Test_API_slice(t *testing.T) {
	type curl struct {
		Header []string `clop:"-H; --header; " 
						 usage:"Pass custom header LINE to server (H)"`
	}

	type curl2 struct {
		Header []string `clop:"-H; --header; greedy" 
						 usage:"Pass custom header LINE to server (H)"`
		Url   string   `clop:"--url" usage:"URL to work with"`
		Count []string `clop:"-c; greedy" usage:"test count"`
	}

	for index, test := range []testAPI{
		// 长选项
		{
			func() curl {
				c := curl{}
				p := New([]string{"--header", "h1:v1", "--header", "h2:v2", "--header", "h3:v3"})
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl{Header: []string{"h1:v1", "h2:v2", "h3:v3"}},
		},
		// 短选项
		{
			func() curl {
				c := curl{}
				p := New([]string{"-H", "h1:v1", "-H", "h2:v2", "-H", "h3:v3"})
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl{Header: []string{"h1:v1", "h2:v2", "h3:v3"}},
		},
		// 长选项 + 贪婪模式
		{
			func() curl2 {
				c := curl2{}
				p := New([]string{"--header", "h1:v1", "h2:v2", "h3:v3", "--url", "www.baidu.com"})
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl2{Header: []string{"h1:v1", "h2:v2", "h3:v3"}, Url: "www.baidu.com"},
		},
		// 短选项+贪婪模式
		{
			func() curl2 {
				c := curl2{}
				p := New([]string{"-H", "h1:v1", "h2:v2", "h3:v3", "--url", "www.baidu.com", "-cval1", "val2", "val3"})
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl2{Header: []string{"h1:v1", "h2:v2", "h3:v3"}, Url: "www.baidu.com", Count: []string{"val1", "val2", "val3"}},
		},
	} {
		assert.Equal(t, test.got, test.need, fmt.Sprintf("test index = %d", index))
	}
}

// 测试int类型
func Test_API_int(t *testing.T) {
	type grep struct {
		BeforeContext int      `clop:"-B;--before-context" usage:"print NUM lines of leading context"`
		AfterContext  int      `clop:"-A;--after-context"   usage:"print NUM lines of trailing context"`
		MaxCount      int      `clop:"-m; --max-count" usage:"Stop reading the file after num matches"`
		Args          []string `clop:"args"`
	}

	for _, test := range []testAPI{
		// 测试短选项
		{
			func() grep {
				g := grep{}
				cp := New([]string{"-B3", "--after-context", "1", "keyword", "join.txt", "-m", "4"})
				err := cp.Bind(&g)
				assert.NoError(t, err)
				return g
			}(), grep{BeforeContext: 3, AfterContext: 1, MaxCount: 4, Args: []string{"keyword", "join.txt"}},
		},
	} {
		assert.Equal(t, test.need, test.got)
	}
}

// 测试环境变量
func Test_API_env(t *testing.T) {
	type env struct {
		Url     []string `clop:"-u; --url; env=CLOP-TEST-URL" usage:"URL to work with"`
		Debug   bool     `clop:"-d; --debug; env=CLOP-DEBUG" usage:"debug"`
		MaxLine int      `clop:"env=CLOP-MAXLINE" usage:"test int"`
	}

	for _, test := range []testAPI{
		{
			func() env {
				e := env{}
				defer func() {
					os.Unsetenv("CLOP-TEST-URL")
					os.Unsetenv("CLOP-DEBUG")
				}()
				os.Setenv("CLOP-TEST-URL", "godoc.org")
				err := os.Setenv("CLOP-DEBUG", "")

				assert.NoError(t, err)

				p := New([]string{"-u", "qq.com", "-u", "baidu.com"})
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{Url: []string{"qq.com", "baidu.com", "godoc.org"}, Debug: true},
		},
		{
			func() env {
				defer func() {
					os.Unsetenv("CLOP-MAXLINE")
				}()
				err := os.Setenv("CLOP-MAXLINE", "3")
				assert.NoError(t, err)

				e := env{}
				p := New([]string{})
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{MaxLine: 3},
		},
		{
			func() env {
				defer func() {
					os.Unsetenv("CLOP-DEBUG")
				}()
				err := os.Setenv("CLOP-DEBUG", "false")
				assert.NoError(t, err)

				e := env{}
				p := New([]string{})
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{},
		},
	} {
		assert.Equal(t, test.need, test.got)
	}
}

// 多行usage消息
func Test_API_head(t *testing.T) {
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

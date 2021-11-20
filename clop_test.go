package clop

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

		Files []string `clop:"args=files"`
	}

	for index, test := range []testAPI{
		// 测试短选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-vTsnEb"}).SetExit(false)
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},
		// 测试短选项,就一个选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-v"}).SetExit(false)
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{ShowNonprinting: true},
		},
		// 测试长选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"--show-nonprinting", "--show-tabs", "--squeeze-blank", "--number", "--show-ends", "--number-nonblank"}).SetExit(false)
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},

		// 测试长短选项混合的情况
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-v", "--show-tabs", "-s", "--number", "-E", "--number-nonblank"}).SetExit(false)
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{NumberNonblank: true, ShowEnds: true, Number: true, SqueezeBlank: true, ShowTab: true, ShowNonprinting: true},
		},

		// 测试args选项
		{
			func() cat {
				c := cat{}
				cp := New([]string{"-n", "r.go", "-T", "pool.c"}).SetExit(false)
				err := cp.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), cat{Number: true, ShowTab: true, Files: []string{"r.go", "pool.c"}},
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
				p := New([]string{"--header", "h1:v1", "--header", "h2:v2", "--header", "h3:v3"}).SetExit(false)
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl{Header: []string{"h1:v1", "h2:v2", "h3:v3"}},
		},
		// 短选项
		{
			func() curl {
				c := curl{}
				p := New([]string{"-H", "h1:v1", "-H", "h2:v2", "-H", "h3:v3"}).SetExit(false)
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl{Header: []string{"h1:v1", "h2:v2", "h3:v3"}},
		},
		// 长选项 + 贪婪模式
		{
			func() curl2 {
				c := curl2{}
				p := New([]string{"--header", "h1:v1", "h2:v2", "h3:v3", "--url", "www.baidu.com"}).SetExit(false)
				err := p.Bind(&c)
				assert.NoError(t, err)
				return c
			}(), curl2{Header: []string{"h1:v1", "h2:v2", "h3:v3"}, Url: "www.baidu.com"},
		},
		// 短选项 + 贪婪模式
		{
			func() curl2 {
				c := curl2{}
				p := New([]string{"-H", "h1:v1", "h2:v2", "h3:v3", "--url", "www.baidu.com", "-cval1", "val2", "val3"}).SetExit(false)
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
		Files         []string `clop:"args=files"`
	}

	for _, test := range []testAPI{
		// 测试短选项
		{
			func() grep {
				g := grep{}
				cp := New([]string{"-B3", "--after-context", "1", "keyword", "join.txt", "-m", "4"}).SetExit(false)
				err := cp.Bind(&g)
				assert.NoError(t, err)
				return g
			}(), grep{BeforeContext: 3, AfterContext: 1, MaxCount: 4, Files: []string{"keyword", "join.txt"}},
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

				p := New([]string{"-u", "qq.com", "-u", "baidu.com"}).SetExit(false)
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
				p := New([]string{}).SetExit(false)
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
				p := New([]string{}).SetExit(false)
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{},
		},
	} {
		assert.Equal(t, test.need, test.got)
	}
}

// args
func Test_API_args(t *testing.T) {
	type testArgs struct {
		Debug  bool     `clop:"-d; --debug" usage:"open debug mode"`
		Level  string   `usage:"log level"`
		Input  string   `clop:"args=input"`
		Format string   `clop:"env=CLOP-FORMAT"`
		Files  []string `clop:"args=files" usage:"files to open"`
	}

	for _, test := range []testAPI{
		// 多个args参数
		{
			func() testArgs {
				a := testArgs{}
				defer func() {
					os.Unsetenv("CLOP-FORMAT")
				}()

				err := os.Setenv("CLOP-FORMAT", "mp3")

				assert.NoError(t, err)

				p := New([]string{"-d", "--level", "info", "output.file", "a.txt", "b.txt"}).SetExit(false)
				err = p.Bind(&a)
				assert.NoError(t, err)
				return a
			}(), testArgs{Debug: true, Level: "info", Input: "output.file", Format: "mp3", Files: []string{"a.txt", "b.txt"}},
		},
	} {

		assert.Equal(t, test.need, test.got)
	}
}

func Test_API_versionAndAbout(t *testing.T) {
	type testVersionAndAbout struct {
		V     string `clop:"version=v0.0.1"`
		About string `clop:"about=a quick start example"`
	}

	type testVersionAndAbout2 struct {
		V     bool `clop:"version=v0.0.1"`
		About bool `clop:"about=a quick start example"`
	}

	for range []error{
		func() error {
			va := testVersionAndAbout{}

			p := New([]string{"-h"}).SetExit(false)

			err := p.Bind(&va)

			assert.NoError(t, err)
			if err != nil {
				return err
			}
			va.V = p.version
			va.About = p.about
			assert.Equal(t, va, testVersionAndAbout{V: "v0.0.1", About: "a quick start example"})
			return nil
		}(),

		func() error {
			va := testVersionAndAbout2{}
			p := New([]string{"-h"}).SetExit(false)
			err := p.Bind(&va)
			assert.NoError(t, err)
			if err != nil {
				return err
			}
			assert.Equal(t, p.version, "v0.0.1")
			assert.Equal(t, p.about, "a quick start example")
			return nil
		}(),
	} {
	}
}

func Test_API_subcommand(t *testing.T) {
	type add struct {
		All      bool     `clop:"-A; --all" usage:"add changes from all tracked and untracked files"`
		Force    bool     `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
		Pathspec []string `clop:"args=pathspec"`
	}

	type mv struct {
		Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
	}

	type git struct {
		Add add `clop:"subcommand=add" usage:"Add file contents to the index"`
		Mv  mv  `clop:"subcommand=mv" usage:"Move or rename a file, a directory, or a symlink"`
	}

	// 测试正确的情况
	for _, test := range []testAPI{
		{
			// 测试add子命令
			func() git {
				g := git{}
				p := New([]string{"add", "-Af", "a.txt"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				assert.True(t, p.IsSetSubcommand("add"))
				assert.False(t, p.IsSetSubcommand("mv"))
				return g
			}(), git{Add: add{All: true, Force: true, Pathspec: []string{"a.txt"}}}},
		{
			// 测试mv子命令
			func() git {
				g := git{}
				p := New([]string{"mv", "-f"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				assert.False(t, p.IsSetSubcommand("add"))
				assert.True(t, p.IsSetSubcommand("mv"))
				return g
			}(), git{Mv: mv{Force: true}}},
		{
			// 测试-h 输出的Usage
			func() git {
				g := git{}
				p := New([]string{"-h"}).SetExit(false)
				b := &bytes.Buffer{}
				p.w = b
				err := p.Bind(&g)
				assert.NoError(t, err)
				assert.True(t, checkUsage(b))
				os.Stdout.Write(b.Bytes())
				return g
			}(), git{Add: add{}}},
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

		ZeroTerminated byte `clop:"-z;--zero-terminated"  usage:"line delimiter is NUL, not newline"`
	}

	h := head{}
	cp := New([]string{}).SetExit(false)
	err := cp.register(&h)
	assert.NoError(t, err)
}

// 测试错误的情况
func Test_API_fail(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-b, --number-nonblank"
		                     usage:"number nonempty output lines, overrides"`
	}

	for range []struct{}{

		func() struct{} {
			c := cat{}
			cp := New([]string{"-vTsnEb"}).SetExit(false)
			err := cp.Bind(&c)
			assert.Error(t, err)
			return struct{}{}
		}(),
	} {
	}
}

// 设置数据校验
func Test_API_valid(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-b; --number-nonblank" valid:"required"
		                     usage:"number nonempty output lines, overrides"`
	}

	// 数据校验 + 子命令测试
	type add struct {
		All      bool     `clop:"-A; --all" usage:"add changes from all tracked and untracked files" valid:"required"`
		Force    bool     `clop:"-f; --force" usage:"allow adding otherwise ignored files" valid:"required"`
		Pathspec []string `clop:"args=pathspec"`
	}

	type mv struct {
		Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files" valid:"required"`
	}

	type git struct {
		Add add `clop:"subcommand=add" usage:"Add file contents to the index"`
		Mv  mv  `clop:"subcommand=mv" usage:"Move or rename a file, a directory, or a symlink"`
	}

	for range []struct{}{

		func() struct{} {
			g := git{}
			cp := New([]string{"mv", "-f"}).SetExit(false)
			err := cp.Bind(&g)
			assert.NoError(t, err)
			assert.Equal(t, g, git{Mv: mv{Force: true}})
			return struct{}{}
		}(),
		func() struct{} {
			c := cat{}
			cp := New([]string{}).SetExit(false)
			err := cp.Bind(&c)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			g := git{}
			cp := New([]string{"add", "-A", "-f"}).SetExit(false)
			err := cp.Bind(&g)
			assert.NoError(t, err)
			assert.Equal(t, g, git{Add: add{All: true, Force: true}})
			return struct{}{}
		}(),
	} {
	}
}

func Test_API_valid_fail(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-b; --number-nonblank" valid:"xxx"
		                     usage:"number nonempty output lines, overrides"`
	}

	for range []struct{}{

		func() struct{} {
			c := cat{}
			cp := New([]string{}).SetExit(false)
			assert.Panics(t, func() {
				err := cp.Bind(&c)
				assert.NoError(t, err)
			})
			return struct{}{}
		}(),
	} {
	}
}

func Test_Option_checkOptionName(t *testing.T) {
	// 测试错误的情况
	for _, v := range []string{
		"c,--bytes",
		"c --bytes",
	} {
		_, b := checkOptionName(v)
		assert.False(t, b)
	}
	//测试正确的情况
	for _, v := range []string{
		"1",
		"2",
		"c",
		"bytes",
		"number-nonblank",
		"pkg_add",
	} {
		_, b := checkOptionName(v)
		assert.True(t, b, fmt.Sprintf("option name is :%s", v))
	}
}

// 测试重复值报错
func Test_DupTag(t *testing.T) {
	type dup struct {
		Number  int `clop:"-n; --number" usage:"number"`
		Number2 int `clop:"-n" usage:"number"`
	}

	type dup2 struct {
		Number  int `clop:"-n; --number" usage:"number"`
		Number2 int `clop:"--number" usage:"number"`
	}

	for range []struct{}{
		func() struct{} {
			var o bytes.Buffer
			d := dup{}
			p := New([]string{}).SetOutput(&o).SetExit(false)
			err := p.Bind(&d)
			assert.Error(t, err)
			assert.Equal(t, o.String(), "-n is already in use\nFor more information try --help\n")
			return struct{}{}
		}(),
		func() struct{} {
			var o bytes.Buffer
			d := dup2{}
			p := New([]string{}).SetOutput(&o).SetExit(false)
			err := p.Bind(&d)
			assert.Error(t, err)
			assert.Equal(t, o.String(), "--number is already in use\nFor more information try --help\n")
			return struct{}{}
		}(),
	} {
	}
}

// 测试没有注册的选项
func Test_unregistered(t *testing.T) {
	type cat struct {
		Number int `clop:"-n; --number" usage:"number"`
	}

	for range []struct{}{
		func() struct{} {
			var o bytes.Buffer
			c := cat{}
			p := New([]string{"-x"}).SetOutput(&o).SetExit(false)
			err := p.Bind(&c)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			var o bytes.Buffer
			c := cat{}
			p := New([]string{"--www"}).SetOutput(&o).SetExit(false)
			err := p.Bind(&c)
			assert.Error(t, err)
			return struct{}{}
		}(),
	} {
	}
}

// 测试长短选项错误的情况
func Test_shortLongFail(t *testing.T) {
	type shortLong struct {
		Debug bool `clop:"-d; --debug" usage:"debug mode"`
	}

	for _, err := range []error{
		func() error {
			s := shortLong{}
			p := New([]string{"-debug"}).SetExit(false)
			return p.Bind(&s)
		}(),
		func() error {
			s := shortLong{}
			p := New([]string{"--d"}).SetExit(false)
			return p.Bind(&s)
		}(),
	} {

		assert.Error(t, err)
	}
}

// 测试没有clop,只有usage的情况
func Test_noclop(t *testing.T) {
	type Opt struct {
		Debug   bool   `usage:"Activate debug mode"`
		Verbose []bool `usage:"Verbose mode (-v, -vv, -vvv, etc.)"`
	}

	for range []struct{}{
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d", "-v", "-v", "-v"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: true, Verbose: []bool{true, true, true}})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d", "-vvv"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: true, Verbose: []bool{true, true, true}})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d", "--verbose", "--verbose", "--verbose"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: true, Verbose: []bool{true, true, true}})
			return struct{}{}
		}(),
	} {
	}
}

func Test_Curl(t *testing.T) {
	type Curl struct {
		Method string   `clop:"-X; --request" usage:"Specify request command to use"`
		Header []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
		Data   string   `clop:"-d; --data"   usage:"HTTP POST data"`
		Form   []string `clop:"-F; --form"  usage:"Specify multipart MIME data"`
		URL    string   `clop:"args=url" usage:"url"`
	}

	for range []struct{}{
		func() struct{} {
			o := Curl{}
			p := New([]string{"-X", "POST", "-H", "h1:v1", "-H", "h2:v2", "http://127.0.0.1:42397"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Curl{Method: "POST", Header: []string{"h1:v1", "h2:v2"}, URL: "http://127.0.0.1:42397"})
			return struct{}{}
		}(),
		func() struct{} {
			o := Curl{}
			p := New([]string{"-X", "POST", "--header", "h1:v1", "--header", "h2:v2", "http://127.0.0.1:42397"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Curl{Method: "POST", Header: []string{"h1:v1", "h2:v2"}, URL: "http://127.0.0.1:42397"})
			return struct{}{}
		}(),
	} {
	}
}

// 测试选项值中带有=号
func Test_EqualSign(t *testing.T) {
	type Opt struct {
		Debug   bool     `clop:"-d; --debug", usage:"Activate debug mode" defaut:"true"`
		Level   string   `clop:"-l; --level" usage:"log level"`
		Files   []string `clop:"-f; --files" usage:"file"`
		Verbose []bool   `usage:"Verbose mode (-v, -vv, -vvv, etc.)"`
	}

	for range []struct{}{
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "-f=a.txt", "-f=b.txt", "-l=info", "-v=false", "-v=true"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: false, Level: "info", Files: []string{"a.txt", "b.txt"}, Verbose: []bool{false, true}})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"--debug=false"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: false})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"--level=info"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Level: "info"})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--level=info"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: false, Level: "info"})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--level=info", "--files=a.txt", "--files=b.txt", "--files=c.txt"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: false, Level: "info", Files: []string{"a.txt", "b.txt", "c.txt"}})
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--verbose=true", "--files=false", "--files=true"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, Opt{Debug: false, Files: []string{"false", "true"}, Verbose: []bool{true}})
			return struct{}{}
		}(),
	} {
	}
}

func Test_OptionPriority(t *testing.T) {
	type cat struct {
		NumberNonblank bool `clop:"-b;--number-nonblank"
		                     usage:"number nonempty output lines, overrides"`

		ShowEnds bool `clop:"-E;--show-ends"
		               usage:"display $ at end of each line"`
	}

	type Opt struct {
		Debug   bool     `clop:"-d; --debug", usage:"Activate debug mode" defaut:"true"`
		Level   string   `clop:"-l; --level" usage:"log level"`
		Files   []string `clop:"-f; --files" usage:"file"`
		Verbose []bool   `usage:"Verbose mode (-v, -vv, -vvv, etc.)"`
	}

	type curl struct {
		URL  string `clop:"--url" usage:"url"`
		URL2 string `clop:"args=url2" usage:"url2"`
	}

	// 环境变量暂时没有优先级
	type env struct {
		args1 string `clop:"env=test-args1" usage:"env 1"`
		args2 string `clop:"env=test-args2" usage:"env 1"`
	}

	for range []struct{}{
		func() struct{} {
			defer func() {
				os.Unsetenv("test-args1")
				os.Unsetenv("test-args2")
			}()

			os.Setenv("test-args1", "godoc.org")
			os.Setenv("test-args2", "godoc.org2")
			c := env{}
			p := New([]string{}).SetExit(false)
			err := p.Bind(&c)
			assert.NoError(t, err)
			assert.Equal(t, p.GetIndex("test-args1"), p.GetIndex("test-args2"))
			return struct{}{}
		}(),
		func() struct{} {
			c := cat{}
			p := New([]string{"-bE"}).SetExit(false)
			err := p.Bind(&c)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("number-nonblank"), p.GetIndex("show-ends"))
			return struct{}{}
		}(),
		func() struct{} {
			c := cat{}
			p := New([]string{"-Eb"}).SetExit(false)
			err := p.Bind(&c)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("show-ends"), p.GetIndex("number-nonblank"))
			return struct{}{}
		}(),
		func() struct{} {
			c := curl{}
			p := New([]string{"--url", "url", "url2"}).SetExit(false)
			err := p.Bind(&c)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("url"), p.GetIndex("url2"))
			return struct{}{}
		}(),
		func() struct{} {
			c := curl{}
			p := New([]string{"url2", "--url", "url"}).SetExit(false)
			err := p.Bind(&c)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("url2"), p.GetIndex("url"))
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--level=info"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("debug"), p.GetIndex("level"))
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "-f=a.txt", "-f=b.txt", "-l=info", "-v=false", "-v=true"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)

			assert.Less(t, p.GetIndex("debug"), p.GetIndex("files"))
			assert.Less(t, p.GetIndex("files"), p.GetIndex("verbose"))
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--level=info", "--files=a.txt", "--files=b.txt", "--files=c.txt"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("debug"), p.GetIndex("level"))
			assert.Less(t, p.GetIndex("debug"), p.GetIndex("files"))
			assert.Less(t, p.GetIndex("level"), p.GetIndex("files"))
			return struct{}{}
		}(),
		func() struct{} {
			o := Opt{}
			p := New([]string{"-d=false", "--verbose=true", "--files=false", "--files=true"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Less(t, p.GetIndex("debug"), p.GetIndex("verbose"))
			assert.Less(t, p.GetIndex("debug"), p.GetIndex("files"))
			assert.Less(t, p.GetIndex("verbose"), p.GetIndex("files"))
			return struct{}{}
		}(),
	} {
	}
}

func Test_OverloadHelp(t *testing.T) {
	type help struct {
		Help  bool `clop:"-h;--help" usage:"Overload help"`
		Debug bool `clop:"-d;--debug" usage:"debug mode"`
	}

	for range []struct{}{
		func() struct{} {
			h := help{}
			p := New([]string{"-h", "-d"}).SetExit(false)
			err := p.Bind(&h)
			assert.NoError(t, err)
			assert.Equal(t, h, help{Help: true, Debug: true})
			return struct{}{}
		}(),
		func() struct{} {
			h := help{}
			p := New([]string{"--help"}).SetExit(false)
			err := p.Bind(&h)
			assert.NoError(t, err)
			assert.Equal(t, h, help{Help: true})
			return struct{}{}
		}(),
	} {
	}
}

// 测试内部register接口
func Test_Internal_register(t *testing.T) {
	p := New([]string{}).SetExit(false)
	for k, v := range []interface{}{
		nil,
		struct{}{},
		(*int)(nil),
	} {

		err := p.Bind(v)
		assert.Error(t, err, fmt.Sprintf("test case index:%d", k))
	}
}

// 测试一个-号的情况
func Test_One_(t *testing.T) {
	type One struct {
		Debug  bool   `clop:"-d;--debug" usage:"debug mode"`
		Stdout string `clop:"args=stdout" usage:"stdout"`
	}

	for range []struct{}{
		func() struct{} {
			o := One{}
			p := New([]string{"-d", "-"}).SetExit(false)
			err := p.Bind(&o)
			assert.NoError(t, err)
			assert.Equal(t, o, One{Debug: true, Stdout: "-"})
			return struct{}{}
		}(),
	} {
	}
}

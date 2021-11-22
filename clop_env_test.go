package clop

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEnvHelp struct {
	Number   int  `clop:"long;env=NUMBER" usage:"number all output lines"`
	ShowEnds bool `clop:"-E;long;env=SHOW_ENDS" usage:"display $ at end of each line"`
}

type testEnv2Help struct {
	Number   int  `clop:"long;env=NUMBER" usage:"number all output lines"`
	ShowEnds bool `clop:"-E;long;env=SHOW_ENDS" usage:"display $ at end of each line"`
}

// 常规env help测试
func Test_Env_usage(t *testing.T) {
	te := testEnvHelp{}
	var out bytes.Buffer
	c := New([]string{"--help"}).SetExit(false).SetOutput(&out)
	c.Bind(&te)
	assert.NotEqual(t, -1, bytes.Index(out.Bytes(), []byte("Environment Variable:")))
}

// 简写env help测试
func Test_Env2_usage(t *testing.T) {
	te := testEnv2Help{}
	var out bytes.Buffer
	c := New([]string{"--help"}).SetExit(false).SetOutput(&out)
	c.Bind(&te)
	assert.NotEqual(t, -1, bytes.Index(out.Bytes(), []byte("Environment Variable:")))
}

// 测试常规用法环境变量
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

func Test_API_env2(t *testing.T) {
	type env struct {
		ClopTestUrl []string `clop:"-u; --url; env" usage:"URL to work with"`
		ClopDebug   bool     `clop:"-d; --debug; env" usage:"debug"`
		ClopMaxline int      `clop:"env" usage:"test int"`
	}

	for _, test := range []testAPI{
		{
			func() env {
				e := env{}
				defer func() {
					os.Unsetenv("CLOP_TEST_URL")
					os.Unsetenv("CLOP_DEBUG")
				}()
				os.Setenv("CLOP_TEST_URL", "godoc.org")
				err := os.Setenv("CLOP_DEBUG", "")

				assert.NoError(t, err)

				p := New([]string{"-u", "qq.com", "-u", "baidu.com"}).SetExit(false)
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{ClopTestUrl: []string{"qq.com", "baidu.com", "godoc.org"}, ClopDebug: true},
		},
		{
			func() env {
				defer func() {
					os.Unsetenv("CLOP_MAXLINE")
				}()
				err := os.Setenv("CLOP_MAXLINE", "3")
				assert.NoError(t, err)

				e := env{}
				p := New([]string{}).SetExit(false)
				err = p.Bind(&e)
				assert.NoError(t, err)
				return e
			}(), env{ClopMaxline: 3},
		},
		{
			func() env {
				defer func() {
					os.Unsetenv("CLOP_DEBUG")
				}()
				err := os.Setenv("CLOP_DEBUG", "false")
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

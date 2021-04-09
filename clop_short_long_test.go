package clop

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试短选项的情况
func Test_Short(t *testing.T) {
	type short struct {
		Int         int      `clop:"short" default:"1"`
		Float64     float64  `clop:"short" default:"3.64"`
		SliceString []string `clop:"short" default:"[\"one\", \"two\"]"`
		Name        []int    `clop:"short" usage:"slice test" valid:"required" default:"[1,2]"`
	}

	defaultShort := short{
		Int:         1,
		Float64:     3.64,
		SliceString: []string{"one", "two"},
		Name:        []int{1, 2},
	}

	for range []struct{}{
		// 正常用法
		func() struct{} {
			got := short{}
			p := New([]string{"-i", "333", "-f", "4444", "-s", "3", "-s", "4", "-n", "3", "-n", "4"}).SetExit(false)
			err := p.Bind(&got)

			need := short{
				Int:         333,
				Float64:     4444.0,
				SliceString: []string{"3", "4"},
				Name:        []int{3, 4},
			}

			assert.Equal(t, need, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值, 没有命令行选项的情况
		func() struct{} {
			got := short{}
			p := New([]string{}).SetExit(false)
			err := p.Bind(&got)
			assert.Equal(t, defaultShort, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值，有些值没有命令行选项的，有些使用默认值
		func() struct{} {
			got := short{}
			need := short{
				Int:         333,
				Float64:     3.64,
				SliceString: []string{"one", "two"},
				Name:        []int{3, 4},
			}

			p := New([]string{"-i", "333", "-n", "3", "-n", "4"}).SetExit(false)
			p.Bind(&got)
			assert.Equal(t, got, need)
			return struct{}{}
		}(),

		// 测试-h帮助信息
		func() struct{} {
			got := short{}
			var out bytes.Buffer
			p := New([]string{"-h"}).SetExit(false).SetOutput(&out)
			err := p.Bind(&got)
			assert.NoError(t, err)

			needTest := []string{
				`1`,
				`3.64`,
				`["one", "two"]`,
				`[1,2]`,
			}

			helpMessage := out.String()
			for _, v := range needTest {
				pos := strings.Index(helpMessage, v)
				assert.NotEqual(t, pos, -1, fmt.Sprintf("search (%s) not found", v))
			}
			return struct{}{}
		}(),
	} {
	}
}

func Test_Long(t *testing.T) {
	type short struct {
		Int         int      `clop:"long" default:"1"`
		Float64     float64  `clop:"long" default:"3.64"`
		SliceString []string `clop:"long" default:"[\"one\", \"two\"]"`
		Name        []int    `clop:"long" usage:"slice test" valid:"required" default:"[1,2]"`
	}

	defaultShort := short{
		Int:         1,
		Float64:     3.64,
		SliceString: []string{"one", "two"},
		Name:        []int{1, 2},
	}

	for range []struct{}{
		// 正常用法
		func() struct{} {
			got := short{}
			p := New([]string{"--int", "333", "--float64", "4444", "--slice-string", "3", "--slice-string", "4", "--name", "3", "--name", "4"}).SetExit(false)
			err := p.Bind(&got)

			need := short{
				Int:         333,
				Float64:     4444.0,
				SliceString: []string{"3", "4"},
				Name:        []int{3, 4},
			}

			assert.Equal(t, need, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值, 没有命令行选项的情况
		func() struct{} {
			got := short{}
			p := New([]string{}).SetExit(false)
			err := p.Bind(&got)
			assert.Equal(t, defaultShort, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值，有些值没有命令行选项的，有些使用默认值
		func() struct{} {
			got := short{}
			need := short{
				Int:         333,
				Float64:     3.64,
				SliceString: []string{"one", "two"},
				Name:        []int{3, 4},
			}

			p := New([]string{"--int", "333", "--name", "3", "--name", "4"}).SetExit(false)
			p.Bind(&got)
			assert.Equal(t, got, need)
			return struct{}{}
		}(),

		// 测试-h帮助信息
		func() struct{} {
			got := short{}
			var out bytes.Buffer
			p := New([]string{"-h"}).SetExit(false).SetOutput(&out)
			err := p.Bind(&got)
			assert.NoError(t, err)

			needTest := []string{
				`1`,
				`3.64`,
				`["one", "two"]`,
				`[1,2]`,
			}

			helpMessage := out.String()
			for _, v := range needTest {
				pos := strings.Index(helpMessage, v)
				assert.NotEqual(t, pos, -1, fmt.Sprintf("search (%s) not found", v))
			}
			return struct{}{}
		}(),
	} {
	}
}

func Test_LongAndShort(t *testing.T) {
	type short struct {
		Int         int      `clop:"long;short" default:"1"`
		Float64     float64  `clop:"long;short" default:"3.64"`
		SliceString []string `clop:"long;short" default:"[\"one\", \"two\"]"`
		Name        []int    `clop:"long;short" usage:"slice test" valid:"required" default:"[1,2]"`
	}

	defaultShort := short{
		Int:         1,
		Float64:     3.64,
		SliceString: []string{"one", "two"},
		Name:        []int{1, 2},
	}

	for range []struct{}{
		// 正常用法
		func() struct{} {
			got := short{}
			p := New([]string{"--int", "333", "--float64", "4444", "--slice-string", "3", "-s", "4", "-n", "3", "--name", "4"}).SetExit(false)
			err := p.Bind(&got)

			need := short{
				Int:         333,
				Float64:     4444.0,
				SliceString: []string{"3", "4"},
				Name:        []int{3, 4},
			}

			assert.Equal(t, need, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值, 没有命令行选项的情况
		func() struct{} {
			got := short{}
			p := New([]string{}).SetExit(false)
			err := p.Bind(&got)
			assert.Equal(t, defaultShort, got)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		// 测试默认值，有些值没有命令行选项的，有些使用默认值
		func() struct{} {
			got := short{}
			need := short{
				Int:         333,
				Float64:     3.64,
				SliceString: []string{"one", "two"},
				Name:        []int{3, 4},
			}

			p := New([]string{"-i", "333", "--name", "3", "-n", "4"}).SetExit(false)
			p.Bind(&got)
			assert.Equal(t, got, need)
			return struct{}{}
		}(),

		// 测试-h帮助信息
		func() struct{} {
			got := short{}
			var out bytes.Buffer
			p := New([]string{"-h"}).SetExit(false).SetOutput(&out)
			err := p.Bind(&got)
			assert.NoError(t, err)

			needTest := []string{
				`1`,
				`3.64`,
				`["one", "two"]`,
				`[1,2]`,
			}

			helpMessage := out.String()
			for _, v := range needTest {
				pos := strings.Index(helpMessage, v)
				assert.NotEqual(t, pos, -1, fmt.Sprintf("search (%s) not found", v))
			}
			return struct{}{}
		}(),
	} {
	}
}

func Test_Error_Usage(t *testing.T) {
	type short struct {
		Int int `clop:"long;short" valid:"required"`
	}

	var out bytes.Buffer
	p := New([]string{}).SetExit(false).SetOutput(&out)
	p.Bind(&short{})

	needTest := []string{
		"-i",
		"--int",
	}

	usage := out.Bytes()
	for _, n := range needTest {
		pos := bytes.Index(usage, []byte(n))
		assert.NotEqual(t, pos, -1, fmt.Sprintf("search (%s) not found", n))
	}
	//io.Copy(os.Stdout, &out)
}

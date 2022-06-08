package clop

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type add struct {
	All      bool     `clop:"-A; --all" usage:"add changes from all tracked and untracked files"`
	Force    bool     `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
	Pathspec []string `clop:"args=pathspec"`
	isSet    bool
}

func (a *add) SubMain() {
	a.isSet = true
}

type mv struct {
	Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
	isSet bool
}

func (m *mv) SubMain() {
	m.isSet = true
}

type touch struct {
	Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
	isSet bool
}

func (t *touch) SubMain() {
	t.isSet = true
}

type git struct {
	Add   add   `clop:"subcommand=add" usage:"Add file contents to the index"`
	Mv    mv    `clop:"subcommand=mv" usage:"Move or rename a file, a directory, or a symlink"`
	Touch touch `clop:"subcommand" usage:"touch file to the index"`
}

// 方法1
// 子命令自动调用SubMain接口
func Test_API_subcommand_SubMain(t *testing.T) {

	// 测试正确的情况
	for _, test := range []testAPI{
		{
			// 测试add子命令
			func() git {
				g := git{}
				p := New([]string{"add", "-Af", "a.txt"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				assert.True(t, g.Add.isSet)
				assert.False(t, g.Mv.isSet)
				return g
			}(), git{Add: add{All: true, Force: true, isSet: true, Pathspec: []string{"a.txt"}}}},
		{
			// 测试mv子命令
			func() git {
				g := git{}
				p := New([]string{"mv", "-f"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				return g
			}(), git{Mv: mv{Force: true, isSet: true}}},
		{
			// 测试touch子命令
			func() git {
				g := git{}
				p := New([]string{"touch", "-f"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				return g
			}(), git{Touch: touch{Force: true, isSet: true}}},
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

// 方法2
func Test_API_subcommand(t *testing.T) {

	type add struct {
		All      bool     `clop:"-A; --all" usage:"add changes from all tracked and untracked files"`
		Force    bool     `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
		Pathspec []string `clop:"args=pathspec"`
		isSet    bool
	}

	type mv struct {
		Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
		isSet bool
	}

	type touch struct {
		Force bool `clop:"-f; --force" usage:"allow adding otherwise ignored files"`
		isSet bool
	}

	type git struct {
		Add   add   `clop:"subcommand=add" usage:"Add file contents to the index"`
		Mv    mv    `clop:"subcommand=mv" usage:"Move or rename a file, a directory, or a symlink"`
		Touch touch `clop:"subcommand" usage:"touch file to the index"`
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
			// 测试touch子命令
			func() git {
				g := git{}
				p := New([]string{"touch", "-f"}).SetExit(false)
				err := p.Bind(&g)
				assert.NoError(t, err)
				assert.False(t, p.IsSetSubcommand("add"))
				assert.True(t, p.IsSetSubcommand("touch"))
				return g
			}(), git{Touch: touch{Force: true}}},
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

package clop

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试显示信息, 普通命令
func Test_Version_Show2(t *testing.T) {
	type show struct {
		Max int `clop:"short;long" usage:"max threads"`
	}

	var buf bytes.Buffer
	cmd := New([]string{"-h"}).SetProcName("Test_Version_Show").SetVersion("v0.0.2").SetExit(false).SetOutput(&buf)
	cmd.MustBind(&show{})

	assert.NotEqual(t, strings.Index(buf.String(), "v0.0.2"), -1)
	assert.NotEqual(t, strings.Index(buf.String(), "print version information"), -1)
	assert.NotEqual(t, strings.Index(buf.String(), "print the help information"), -1)
}

// 测试显示信息
func Test_Version_Show(t *testing.T) {
	var buf bytes.Buffer
	cmd := New([]string{"-h"}).SetProcName("Test_Version_Show").SetVersion("v0.0.2").SetExit(false).SetOutput(&buf)
	cmd.MustBind(&git{})

	assert.NotEqual(t, strings.Index(buf.String(), "v0.0.2"), -1)
}

// 测试about信息
func Test_About_Show(t *testing.T) {

	var buf bytes.Buffer
	cmd := New([]string{"-h"}).SetProcName("Test_Version_Show").SetAbout("about信息").SetExit(false).SetOutput(&buf)
	cmd.MustBind(&git{})

	assert.NotEqual(t, strings.Index(buf.String(), "about信息"), -1)
}

// 测试-V
func Test_Version_Option_Short(t *testing.T) {
	var buf bytes.Buffer
	procName := "Test_Version_Option_Short"
	version := "v0.2.0"
	cmd := New([]string{"-V"}).SetProcName(procName).SetVersion(version).SetExit(false).SetOutput(&buf)
	cmd.MustBind(&git{})

	assert.NotEqual(t, strings.Index(buf.String(), fmt.Sprintf("%s %s\n", procName, version)), -1)
}

// 测试 --version
func Test_Version_Option_Long(t *testing.T) {
	var buf bytes.Buffer
	procName := "Test_Version_Option_Short"
	version := "v0.2.0"
	cmd := New([]string{"--version"}).SetProcName(procName).SetVersion(version).SetExit(false).SetOutput(&buf)
	cmd.MustBind(&git{})

	assert.NotEqual(t, strings.Index(buf.String(), fmt.Sprintf("%s %s\n", procName, version)), -1)
}

// 测试-V
func Test_Version_Option_Short_Replace(t *testing.T) {
	type dup struct {
		Version string `clop:"-V" usage:"usage"`
	}

	d := &dup{}

	var buf bytes.Buffer
	procName := "Test_Version_Option_Short"
	version := "v0.2.0"
	cmd := New([]string{"-V", "1"}).SetProcName(procName).SetVersion(version).SetExit(false).SetOutput(&buf)
	cmd.MustBind(d)

	assert.Equal(t, d.Version, "1")
	assert.Equal(t, strings.Index(buf.String(), fmt.Sprintf("%s %s\n", procName, version)), -1)
}

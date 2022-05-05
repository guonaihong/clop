package clop

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://github.com/guonaihong/clop/issues/86
func Test_Fix86(t *testing.T) {

	type Args struct {
		Server Server `clop:"subcommand=server" usage:"Run in server model"`
	}

	type Server struct {
	}

	var args Args
	var buf bytes.Buffer
	New([]string{"-server"}).SetExit(false).SetOutput(&buf).Bind(&args)
	// 测试是否有panic
	assert.Equal(t, strings.Index(buf.String(), "panic"), -1)
}

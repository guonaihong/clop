package clop

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Str  string   `clop:"--str" valid:"required"`
	List []string `clop:"-l;--list;greedy"`
}

// https://github.com/guonaihong/clop/issues/77
func TestGreedyBug(t *testing.T) {
	need := Config{
		Str:  "aaa",
		List: []string{"bbb", "ccc"},
	}

	for _, v := range [][]string{
		[]string{"--list", "bbb", "ccc", "--str=aaa"},
		[]string{"--list", "bbb", "ccc", "--str", "aaa"},
		[]string{"-l", "bbb", "ccc", "--str=aaa"},
		[]string{"--list", "bbb", "ccc", "--str", "aaa"},
		[]string{"-l", "bbb", "ccc", "--str", "aaa"},
		[]string{"--str=aaa", "--list", "bbb", "ccc"},
		[]string{"--str=aaa", "-l", "bbb", "ccc"},
	} {
		p := New(v).SetExit(false).SetOutput(os.Stdout)
		got := Config{}

		err := p.Bind(&got)
		assert.NoError(t, err)
		assert.Equal(t, need, got)
	}
}

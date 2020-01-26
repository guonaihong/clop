package clop

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type getData struct {
	set  string
	need string
	key  string
}

func Test_Get(t *testing.T) {
	gd := []getData{
		{`clop:"-c;--number-nonblank"
								usage:"number nonempty output lines, overrides"`,
			`-c;--number-nonblank`,
			"clop",
		},

		{`clop:"-c;--number-nonblank"
								usage:"number nonempty output lines, overrides"`,
			`number nonempty output lines, overrides`,
			"usage",
		},
		{`clop:"-s;--squeeze-repeats"
			 usage:"replace each sequence of a repeated character
                            that is listed in the last specified SET,
                            with a single occurrence of that character"`,
			`replace each sequence of a repeated character
                            that is listed in the last specified SET,
                            with a single occurrence of that character`,
			"usage"},
	}

	for _, v := range gd {
		val := Tag(v.set).Get(v.key)
		//fmt.Printf("need(%s) get(%s)\n", v.need, val)
		assert.Equal(t, val, v.need)
	}
}

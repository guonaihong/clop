package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MustBind(t *testing.T) {
	type test struct {
		X int `clop:"short;long" usage:"x" valid:"required"`
	}

	test2 := test{}

	c := New([]string{}).SetExit(false)
	assert.Panics(t, func() {
		c.MustBind(&test2)
	})

}

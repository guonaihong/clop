package clop

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Issue(t *testing.T) {

	type test struct {
		Long int `clop:"short;long" valid:"required"`
	}

	for _, tc := range []string{
		func() string {
			var out bytes.Buffer
			testVal := test{}
			cp := New([]string{"-l"}).SetOutput(&out).SetExit(false)
			err := cp.Bind(&testVal)
			assert.Error(t, err)
			return out.String()
		}(),
		func() string {
			var out bytes.Buffer
			testVal := test{}
			cp := New([]string{"--long"}).SetOutput(&out).SetExit(false)
			err := cp.Bind(&testVal)
			assert.Error(t, err)
			return out.String()
		}(),
	} {

		assert.NotEqual(t, strings.Index(tc, "must have a value!"), -1)
	}
}

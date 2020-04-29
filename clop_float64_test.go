package clop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_API_Float64(t *testing.T) {
	type f struct {
		Start float64 `clop:"--start" usage:"start"`
		End   float64 `clop:"--end" usage:"end"`
	}

	for _, test := range []testAPI{
		{
			func() f {
				f0 := f{}
				cp := New([]string{"--start", "-3", "--end", "-4"}).SetExit(false)
				err := cp.Bind(&f0)
				assert.NoError(t, err)
				return f0
			}(), f{Start: -3, End: -4},
		},
	} {

		assert.Equal(t, test.need, test.got)
	}
}

package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_API_Neg(t *testing.T) {
	type f struct {
		Start     float64 `clop:"-b; --start" usage:"start"`
		End       float64 `clop:"-e; --end" usage:"end"`
		Height    int     `clop:"-h; --height" usage:"height"`
		Separator int8    `clop:"-s; --separator" usage:"separator"`
	}

	for _, test := range []testAPI{
		{
			func() f {
				f0 := f{}
				cp := New([]string{"--start", "-3", "--end", "-4", "--height", "-33", "--separator", "-1"}).SetExit(false)
				err := cp.Bind(&f0)
				assert.NoError(t, err)
				return f0
			}(), f{Start: -3, End: -4, Height: -33, Separator: -1},
		},
		{
			func() f {
				f0 := f{}
				cp := New([]string{"-b", "-3", "-e", "-4", "-h", "-33", "-s", "-1"}).SetExit(false)
				err := cp.Bind(&f0)
				assert.NoError(t, err)
				return f0
			}(), f{Start: -3, End: -4, Height: -33, Separator: -1},
		},
	} {

		assert.Equal(t, test.need, test.got)
	}
}

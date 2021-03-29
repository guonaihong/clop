package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nameTestCase struct {
	in  string
	got string
}

func Test_name(t *testing.T) {
	for _, tcase := range []nameTestCase{
		{"longOpt", "long-opt"},
		{"almost-all", "almost-all"},
		{"almost_all", "almost-all"},
		{"_almost_all", "almost-all"},
		{"LongOpt_all", "long-opt-all"},
	} {
		got, err := gnuOptionName(tcase.in)
		assert.NoError(t, err)
		assert.Equal(t, got, tcase.got)
	}
}

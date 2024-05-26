package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nameTestCase struct {
	in  string
	got string
}

func Test_ShortLongName(t *testing.T) {
	for _, tcase := range []nameTestCase{
		{"longOpt", "long-opt"},
		{"almost-all", "almost-all"},
		{"almost_all", "almost-all"},
		{"_almost_all", "almost-all"},
		{"LongOpt_all", "long-opt-all"},
		{"JSON", "JSON"},
		{"URL", "URL"},
		{"URI", "URI"},
	} {
		got, err := gnuOptionName(tcase.in)
		assert.NoError(t, err)
		assert.Equal(t, got, tcase.got)
	}
}

func Test_EnvName(t *testing.T) {
	for _, tcase := range []nameTestCase{
		{"longOpt", "LONG_OPT"},
		{"almost-all", "ALMOST_ALL"},
		{"almost_all", "ALMOST_ALL"},
		{"_almost_all", "ALMOST_ALL"},
		{"envOpt_all", "ENV_OPT_ALL"},
		{"JSON", "JSON"},
		{"URL", "URL"},
		{"URI", "URI"},
	} {
		got, err := envOptionName(tcase.in)
		assert.NoError(t, err)
		assert.Equal(t, got, tcase.got)
	}
}

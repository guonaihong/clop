package clop

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEnv struct {
	Number   int  `clop:"long;env=NUMBER" usage:"number all output lines"`
	ShowEnds bool `clop:"-E;long;env=SHOW_ENDS" usage:"display $ at end of each line"`
}

func Test_Env_usage(t *testing.T) {
	te := testEnv{}
	var out bytes.Buffer
	c := New([]string{"--help"}).SetExit(false).SetOutput(&out)
	c.Bind(&te)
	assert.NotEqual(t, -1, bytes.Index(out.Bytes(), []byte("Environment Variable:")))
}

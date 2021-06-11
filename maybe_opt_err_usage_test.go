package clop

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Issue7_Test struct {
	Rate string `clop:"short;long"`
	Test bool   `clop:"long"`
}

func Test_Issue64(t *testing.T) {
	var out bytes.Buffer

	got := Issue7_Test{}
	p := New([]string{"-test"}).SetExit(false).SetOutput(&out)
	err := p.Bind(&got)
	assert.Error(t, err)
	assert.NotEqual(t, bytes.Index(out.Bytes(), []byte("test")), -1)
}

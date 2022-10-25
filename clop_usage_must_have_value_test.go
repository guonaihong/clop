package clop

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 该文件定位
// 各种出错的提示信息自测函数

func Test_Usage_MustHaveValue(t *testing.T) {

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

func Test_Usage_UnknownCommand(t *testing.T) {
	type cmd struct {
		C string `clop:"C" usage:""`
	}

	c0 := cmd{}
	var b bytes.Buffer
	c := New([]string{}).SetOutput(&b).SetExit(false)
	c.Bind(&c0)
	assert.NotEqual(t, strings.Index(b.String(), `clop:"short;long"`), -1)
}

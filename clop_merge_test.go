package clop

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试用
type Asr struct {
	ThreadNum int  `clop:"long" usage:"thread number" valid:"required"`
	OpenVad   bool `clop:"long" usage:"open vad" valid:"required"`
}

// 测试用
type Server struct {
	ServerAddress string        `clop:"long" usage:"Server address" valid:"required"`
	Rate          time.Duration `clop:"long" usage:"The speed at which audio is sent" valid:"required"`
}

// 1.测试多结构体串联的help功能
func Test_Merge_Help(t *testing.T) {
	type TestA struct {
		Aa string `clop:"long" usage:"a"`
		Bb string `clop:"long" usage:"b"`
	}

	type TestB struct {
		Cc string `clop:"long" usage:"c"`
		Dd string `clop:"long" usage:"c"`
	}

	var out bytes.Buffer

	p := New([]string{"-h"}).SetExit(false).SetOutput(&out)

	p.Register(&TestA{})

	p.Bind(&TestB{})

	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte("aa")))
	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte("bb")))
	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte("cc")))
	assert.Equal(t, 1, bytes.Count(out.Bytes(), []byte("dd")))

}

// 2.测试多结构体串联的parse功能
func Test_Merge_Parse(t *testing.T) {
	p := New([]string{"--server-address", ":8080", "--rate", "1s", "--thread-num", "20", "--open-vad"}).SetExit(false)

	asr := Asr{}
	ser := Server{}
	p.Register(&asr)
	p.Bind(&ser)

	assert.Equal(t, asr.ThreadNum, 20)
	assert.True(t, asr.OpenVad)
	assert.Equal(t, ser.ServerAddress, ":8080")
	assert.Equal(t, ser.Rate, time.Second)

}

// 3.测试MustRegister接口
func Test_Merge_MustRegister(t *testing.T) {
	assert.Panics(t, func() {
		MustRegister(nil)
	})
}

// 4.测试结构体串联的数据校验功能
func Test_Merge_Valid(t *testing.T) {
	var usage bytes.Buffer

	p := New([]string{"--server-address", ":8080", "--rate", "1s"}).SetExit(false).SetOutput(&usage)

	asr := Asr{}
	ser := Server{}
	p.Register(&asr)
	p.Bind(&ser)

	//os.Stdout.Write(usage.Bytes())
	assert.True(t, bytes.Contains(usage.Bytes(), []byte("must have a value")))
}

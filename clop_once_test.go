package clop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//测试once功能打开
func Test_Once_Open(t *testing.T) {
	type once struct {
		Debug bool `clop:"-d; --debug; once" usage:"debug mode"`
	}

	type onceString struct {
		Addr string `clop:"-a; --addr; once" usage:"server address"`
	}

	type onceInt struct {
		MaxThread int `clop:"-t; --max-thread; once" usage:"max thread"`
	}

	for range []struct{}{
		func() struct{} {
			h := once{}
			p := New([]string{"-d", "-d"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := once{}
			p := New([]string{"--debug", "-d"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := once{}
			p := New([]string{"-d", "--debug"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := once{}
			p := New([]string{"--debug", "--debug"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := onceString{}
			p := New([]string{"-a", ":8080", "-a", ":1234"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := onceString{}
			p := New([]string{"--addr", ":8080", "-a", ":1234"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := onceString{}
			p := New([]string{"--a", ":8080", "--addr", ":1234"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := onceString{}
			p := New([]string{"--addr", ":8080", "--addr", ":1234"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			return struct{}{}
		}(),
		func() struct{} {
			h := onceInt{}
			p := New([]string{"--max-thread", "20", "--max-thread", "50"}).SetExit(false)
			err := p.Bind(&h)
			assert.Error(t, err)
			assert.Equal(t, h, onceInt{MaxThread: 20})
			return struct{}{}
		}(),
	} {
	}
}

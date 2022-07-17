package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCallback struct {
	Size int `clop:"short;long;callback=ParseSize" usage:"parse size"`
	Max  int `clop:"short;long"`
}

func (t *TestCallback) ParseSize(val string) {
	t.Size = 1024 * 1024
}

// 指定函数名
func Test_Callback_SpecifyTheFunctionName(t *testing.T) {
	got := TestCallback{}
	need := TestCallback{Size: 1024 * 1024, Max: 10}
	p := New([]string{"--size", "1MB", "--max", "10"}).SetExit(false)
	err := p.Bind(&got)

	assert.Equal(t, got, need)
	assert.NoError(t, err)
}

type TestCallbackDefault struct {
	Size int `clop:"short;long;callback" usage:"parse size"`
}

// 这是默认函数名
func (t *TestCallbackDefault) Parse(val string) {
	t.Size = 1024 * 1024
}

// 指定函数名
func Test_Callback_Default(t *testing.T) {
	got := TestCallback{}
	need := TestCallback{Size: 1024 * 1024, Max: 10}
	p := New([]string{"--size", "1MB", "--max", "10"}).SetExit(false)
	err := p.Bind(&got)

	assert.Equal(t, got, need)
	assert.NoError(t, err)
}

type TestCallbackPanic struct {
	Size int `clop:"short;long;callback" usage:"parse size"`
}

// 这是默认函数名
func (t *TestCallbackPanic) Parse() {
	t.Size = 1024 * 1024
}

func Test_Callback_Panic(t *testing.T) {
	got := TestCallbackPanic{}
	assert.Panics(t, func() {
		p := New([]string{"--size", "1MB", "--max", "10"}).SetExit(false)
		p.Bind(&got)
	})

}

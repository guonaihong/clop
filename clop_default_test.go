package clop

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试默认(default 标签)值
func Test_DefautlValue(t *testing.T) {
	type defaultExample struct {
		Int          int       `clop:"--int" default:"1"`
		Float64      float64   `clop:"--float64" default:"3.64"`
		Float32      float32   `clop:"--float32" default:"3.32"`
		SliceString  []string  `clop:"--slice-string" default:"[\"one\", \"two\"]"`
		SliceInt     []int     `clop:"--slice-int" default:"[1,2,3,4,5]"`
		SliceFloat64 []float64 `clop:"--slice-float64" default:"[1.1,2.2,3.3,4.4,5.5]"`
		Name         []int     `clop:"-e" usage:"slice test" valid:"required" default:"[1,2]"`
	}

	type tool struct {
		Rate string `clop:"-r; --rate" usage:"rate" default:"8000"`
	}

	for range []struct{}{
		func() struct{} {
			got := defaultExample{}
			p := New([]string{"--slice-string", "333", "--slice-string", "4444", "-e", "3", "-e", "4"}).SetExit(false)
			err := p.Bind(&got)

			assert.Equal(t, got.SliceString, []string{"333", "4444"})
			assert.Equal(t, got.Name, []int{3, 4})
			assert.NoError(t, err)
			return struct{}{}
		}(),

		func() struct{} {
			tol := tool{}
			p := New([]string{}).SetExit(false)
			err := p.Bind(&tol)
			assert.NoError(t, err)
			return struct{}{}
		}(),

		func() struct{} {
			got := defaultExample{}
			need := defaultExample{
				Int:          1,
				Float64:      3.64,
				Float32:      3.32,
				SliceString:  []string{"one", "two"},
				SliceInt:     []int{1, 2, 3, 4, 5},
				SliceFloat64: []float64{1.1, 2.2, 3.3, 4.4, 5.5},
				Name:         []int{1, 2},
			}
			p := New([]string{}).SetExit(false)
			p.Bind(&got)
			assert.Equal(t, got, need)
			return struct{}{}
		}(),
		func() struct{} {
			got := defaultExample{}
			var out bytes.Buffer
			p := New([]string{"-h"}).SetExit(false).SetOutput(&out)
			err := p.Bind(&got)
			assert.NoError(t, err)

			needTest := []string{
				`1`,
				`3.64`,
				`3.32`,
				`["one", "two"]`,
				`[1,2,3,4,5]`,
				`[1.1,2.2,3.3,4.4,5.5]`,
			}

			helpMessage := out.String()
			for _, v := range needTest {
				pos := strings.Index(helpMessage, v)
				assert.NotEqual(t, pos, -1, fmt.Sprintf("search (%s) not found", v))
			}
			return struct{}{}
		}(),
	} {
	}
}

package clop

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set_setSlice(t *testing.T) {
	type testSetSlice struct {
		Args []string
	}

	type testSetSliceInt struct {
		Args []int
	}

	for _, v := range []testAPI{
		{
			func() testSetSlice {
				t := &testSetSlice{}
				v := reflect.ValueOf(t)
				setSlice("1", 0, v.Elem().Field(0))
				setSlice("2", 0, v.Elem().Field(0))
				return *t
			}(),
			testSetSlice{Args: []string{"1", "2"}},
		},
		{
			func() testSetSliceInt {
				t := &testSetSliceInt{}
				v := reflect.ValueOf(t)
				setSlice("1", 0, v.Elem().Field(0))
				setSlice("2", 0, v.Elem().Field(0))
				return *t
			}(),
			testSetSliceInt{Args: []int{1, 2}},
		},
	} {

		assert.Equal(t, v.need, v.got)
	}

}

func Test_Set_resetValue(t *testing.T) {

	for _, v := range []interface{}{
		func() interface{} {
			s := "hello"
			return &s
		}(),
		func() interface{} {
			i := 3
			return &i
		}(),
		func() interface{} {
			s := []string{"hello", "word"}
			return &s
		}(),
	} {
		vv := reflect.ValueOf(v).Elem()
		resetValue(vv)
		assert.True(t, vv.IsZero())

	}
}

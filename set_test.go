package clop

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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

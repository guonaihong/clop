package clop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试显示信息
type testShowLongUsage struct {
	tagName string
	in      []string
	need    []string
}

func Test_Valid_showShortLongUsage(t *testing.T) {
	testList := []testShowLongUsage{
		{
			tagName: "Once",
			in:      []string{"short", "long", "short;long"},
			need:    []string{"-o", "--once", "-o;--once"},
		},
		//short + long tag
		{
			tagName: "Header",
			in:      []string{"short", "long", "short;long"},
			need:    []string{"-h", "--header", "-h;--header"},
		},
		//混合
		{
			tagName: "Header",
			in:      []string{"short", "long", "short;long"},
			need:    []string{"-h", "--header", "-h;--header"},
		},
		//混合+带-的变量名
		{
			tagName: "ByteOffset",
			in:      []string{"short;--byte-offset", "-b;long", "short;long"},
			need:    []string{"-b;--byte-offset", "-b;--byte-offset", "-b;--byte-offset"},
		},
		//单字母选项名称
		{
			tagName: "B",
			in:      []string{"short", "-b;long", "short;long"},
			need:    []string{"-b", "-b", "-b"},
		},
		//长选项写在前面，短选项写在后面
		{
			tagName: "ByteOffset",
			in:      []string{"--byte-offset;short", "long;-b", "long;short"},
			need:    []string{"-b;--byte-offset", "-b;--byte-offset", "-b;--byte-offset"},
		},
	}

	for _, test := range testList {
		got := make([]string, len(test.need))
		for i, v := range test.in {
			got[i] = showShortLongUsage(v, test.tagName)
		}

		assert.Equal(t, test.need, got)
	}
}

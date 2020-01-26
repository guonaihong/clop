package clop

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// 长短选项为什么要分开,遍历的数据会更少
type Clop struct {
	short      map[string]*Option
	long       map[string]*Option
	shortRegex []*Option
	longRegex  []*Option
}

type Option struct {
	Name     string      //命令行选项名称
	Pointer  interface{} //存放需要修改的值的地址
	Usage    string      //帮助信息
	DefValue string      //默认值
}

func New(args []string) *Clop {
	return &Clop{}
}

func (c *Clop) parseTag(clop string, usage string) error {
	options := strings.Split(clop, ";")

	for _, o := range options {
		name := ""
		switch {
		case strings.HasPrefix(o, "--"):
			name = o[2:]
		case strings.HasPrefix(o, "-"):
			name = o[1:]
		}

		if len(name) == 0 {
			//TODO return fail
		}

		fmt.Printf("name(%s)\n", name)
	}

	return nil
}

func (c *Clop) registerCore(v reflect.Value, sf reflect.StructField) error {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		clop := Tag(sf.Tag).Get("clop")
		usage := Tag(sf.Tag).Get("usage")
		fmt.Printf("clop(%s), usage(%s)\n", clop, usage)

		c.parseTag(clop, usage)
		return nil
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)
		fmt.Printf("%s\n", sf.Tag)

		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}

		//fmt.Printf("my.index(%d)(1.%s)-->(2.%s)\n", i, Tag(sf.Tag).Get("clop"), Tag(sf.Tag).Get("usage"))
		//fmt.Printf("stdlib.index(%d)(1.%s)-->(2.%s)\n", i, sf.Tag.Get("clop"), sf.Tag.Get("usage"))
		c.registerCore(v.Field(i), sf)
	}

	return nil
}

var emptyField = reflect.StructField{}

func (c *Clop) register(x interface{}) error {
	v := reflect.ValueOf(x)

	if x == nil || v.IsNil() {
		return ErrUnsupportedType
	}

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("%s:got(%T)", ErrNotPointerType, v.Type())
	}

	c.registerCore(v, emptyField)

	return nil
}

func (c *Clop) Bind(x interface{}) error {
	if err := c.register(x); err != nil {
		return err
	}

	return nil
}

func Bind(x interface{}) error {
	return CommandLine.Bind(x)
}

var CommandLine = New(os.Args)

package clop

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"
)

var (
	ErrDuplicateOptions = errors.New("duplicate command options")
	ErrUsageEmpty       = errors.New("usage cannot be empty")
	ErrUnsupported      = errors.New("unsupported command")
	ErrNotFoundName     = errors.New("no command line options found")
)

// 长短选项为什么要分开,遍历的数据会更少
type Clop struct {
	short      map[string]*Option
	long       map[string]*Option
	shortRegex []*Option
	longRegex  []*Option
	args       []string
	saveArgs   reflect.Value
}

type Option struct {
	Pointer      reflect.Value //存放需要修改的值的地址
	Usage        string        //帮助信息
	showDefValue string        //显示默认值
	index        int
}

func New(args []string) *Clop {
	return &Clop{
		short: make(map[string]*Option),
		long:  make(map[string]*Option),
		args:  args,
	}
}

func (c *Clop) setOption(name string, option *Option, m map[string]*Option) error {
	if _, ok := m[name]; ok {
		return fmt.Errorf("%s:%s", ErrDuplicateOptions, name)
	}

	m[name] = option
	return nil
}

func (c *Clop) getOption(arg string, index *int, numMinuses int) error {
	var (
		option     *Option
		shortIndex int
	)

	// 取出option对象
	switch numMinuses {
	case 2: //长选项
		option, _ = c.long[arg]
		if option == nil {
			return fmt.Errorf("not found")
		}
		value := ""
		//TODO确认 posix
		switch option.Pointer.Kind() {
		//bool类型，不考虑false的情况
		case reflect.Bool:
			value = "true"
		default:
			// 如果是长
			if numMinuses == 2 && *index+1 >= len(c.args) {
				return errors.New("wrong long option")
			}

			if numMinuses == 1 {
				value = arg[shortIndex:]
			} else {
				(*index)++
				value = c.args[*index]
			}

		}

		// 赋值
		return setBase(value, option.Pointer)
	case 1: //短选项
		var a rune
		find := false
		for shortIndex, a = range arg {
			//只支持ascii
			if a >= utf8.RuneSelf {
				return errors.New("Illegal character set")
			}

			value := string(byte(a))
			option, _ = c.short[value]
			if option == nil {
				continue
			}

			find = true
			switch option.Pointer.Kind() {
			case reflect.Bool:
				if err := setBase("true", option.Pointer); err != nil {
					return err
				}

			default:
				shortIndex++
				if err := setBase(arg[shortIndex:], option.Pointer); err != nil {
					return err
				}
			}
		}

		if find {
			return nil
		}
	}

	if option == nil {
		return fmt.Errorf("not found")
	}

	return nil
}

func (c *Clop) parseTagAndSetOption(clop string, usage string, v reflect.Value) error {
	options := strings.Split(clop, ";")

	option := &Option{Usage: usage, Pointer: v}

	findName := false
	for _, opt := range options {
		name := ""
		// TODO 检查name的长度
		switch {
		case strings.HasPrefix(opt, "--"):
			name = opt[2:]
			c.setOption(name, option, c.long)
			findName = true
		case strings.HasPrefix(opt, "-"):
			name = opt[1:]
			c.setOption(name, option, c.short)
			findName = true
		case strings.HasPrefix(opt, "def="):
			option.showDefValue = opt[4:]
		default:
			return fmt.Errorf("%s:%s", ErrUnsupported, opt)
		}

		if strings.HasPrefix(opt, "-") && len(name) == 0 {
			return fmt.Errorf("Illegal command line option:%s", opt)
		}

	}

	if !findName {
		return fmt.Errorf("%s:%s", ErrNotFoundName, clop)
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

		if clop == "args" {
			c.saveArgs = v
			return nil
		}

		// clop 可以省略
		if len(clop) == 0 {
			clop = strings.ToLower(sf.Name)
			if len(clop) == 1 {
				clop = "-" + clop
			} else {
				clop = "--" + clop
			}
		}

		// usage  不能为空
		if len(usage) == 0 {
			return ErrUsageEmpty
		}

		c.parseTagAndSetOption(clop, usage, v)
		return nil
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)

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

func (c *Clop) parseOneOption(index *int) error {

	arg := c.args[*index]

	if len(arg) == 0 {
		//TODO return fail
		return errors.New("fail option")
	}

	if arg[0] != '-' {
		setBase(arg, c.saveArgs)
		return nil
	}

	// arg 必须是减号开头的字符串
	numMinuses := 1

	if arg[1] == '-' {
		numMinuses++
	}

	// 暂不支持=号的情况
	// TODO 考虑下要不要加上

	a := arg[numMinuses:]
	return c.getOption(a, index, numMinuses)
}

func (c *Clop) bindStruct() error {

	for i := 0; i < len(c.args); i++ {

		if err := c.parseOneOption(&i); err != nil {
			return err
		}

	}
	return nil
}

func (c *Clop) Bind(x interface{}) error {
	if err := c.register(x); err != nil {
		return err
	}

	return c.bindStruct()
}

func Bind(x interface{}) error {
	return CommandLine.Bind(x)
}

var CommandLine = New(os.Args)

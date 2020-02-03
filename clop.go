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

// 长短选项分为short和long，优点遍历的数据会更少速度更快
type Clop struct {
	short      map[string]*Option
	long       map[string]*Option
	env        map[string]*Option
	shortRegex []*Option //TODO把值用起来
	longRegex  []*Option //TODO把值用起来
	args       []string
	saveArgs   reflect.Value //TODO 测试args为空的情况
}

type Option struct {
	pointer      reflect.Value //存放需要修改的值的reflect.Value类型变量
	usage        string        //帮助信息
	showDefValue string        //显示默认值 TODO把值用起来
	index        int           //表示参数优先级 TODO把值用起来
	showShort    []string      //help显示的短选项
	showLong     []string      //help显示的长选项
	envName      string        //环境变量
	greedy       bool          //贪婪模式 -H a b c 等于-H a -H b -H c
	haveSetEnv   bool          //判断是否设置过环境变量
}

func New(args []string) *Clop {
	return &Clop{
		short: make(map[string]*Option),
		long:  make(map[string]*Option),
		env:   make(map[string]*Option),
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

// 解析长选项
func (c *Clop) parseLong(arg string, index *int) error {
	var option *Option
	option, _ = c.long[arg]
	if option == nil {
		return fmt.Errorf("not found")
	}

	value := ""
	//TODO确认 posix
	switch option.pointer.Kind() {
	//bool类型，不考虑false的情况
	case reflect.Bool:
		value = "true"
	default:
		// 如果是长选项
		if *index+1 >= len(c.args) {
			return errors.New("wrong long option")
		}

		for {

			(*index)++
			if *index >= len(c.args) {
				return nil
			}

			value = c.args[*index]
			// 如果打开贪婪模式，直到遇到-或者最后一个字符才结束
			if strings.HasPrefix(value, "-") {
				(*index)-- //回退这个选项
				return nil
			}

			if err := setBase(value, option.pointer); err != nil {
				return err
			}

			if option.pointer.Kind() != reflect.Slice && !option.greedy {
				return nil
			}
		}
	}
	return setBase(value, option.pointer)
}

// 设置环境变量
func (o *Option) setEnv() (err error) {
	if o.haveSetEnv {
		return nil
	}

	if len(o.envName) > 0 {
		if v, ok := os.LookupEnv(o.envName); ok {
			if o.pointer.Kind() == reflect.Bool {
				if v != "false" {
					v = "true"
				}
			}

			err := setBase(v, o.pointer)
			if err == nil {
				o.haveSetEnv = true
			}
			return err
		}
	}
	return nil
}

func (c *Clop) parseShort(arg string, index *int) error {
	var (
		option     *Option
		shortIndex int
	)
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
		switch option.pointer.Kind() {
		case reflect.Bool:
			if err := setBase("true", option.pointer); err != nil {
				return err
			}

		default:
			shortIndex++
			for value := arg; ; {

				if len(value[shortIndex:]) > 0 {
					if err := setBase(value[shortIndex:], option.pointer); err != nil {
						return err
					}

					if option.pointer.Kind() != reflect.Slice && !option.greedy {
						return nil
					}
				}

				shortIndex = 0

				(*index)++
				if *index >= len(c.args) {
					return nil
				}
				value = c.args[*index]

				if strings.HasPrefix(value, "-") {
					(*index)--
					return nil
				}

			}

		}
	}

	if find {
		return nil
	}
	return nil
}

func (c *Clop) getOptionAndSet(arg string, index *int, numMinuses int) error {
	// 输出帮助信息
	if arg == "h" || arg == "help" {
		c.printHelpMessage()
		os.Exit(0)
	}
	// 取出option对象
	switch numMinuses {
	case 2: //长选项
		return c.parseLong(arg, index)
	case 1: //短选项
		return c.parseShort(arg, index)
	}

	return nil
}

func (c *Clop) genHelpMessage(h *Help) {

	used := make(map[*Option]struct{}, len(c.short))

	saveHelp := func(options map[string]*Option) {
		for _, v := range options {
			if _, ok := used[v]; ok {
				continue
			}

			used[v] = struct{}{}

			var oneArgs []string

			for _, v := range v.showShort {
				oneArgs = append(oneArgs, "-"+v)
			}

			for _, v := range v.showLong {
				oneArgs = append(oneArgs, "--"+v)
			}

			env := ""
			if len(v.envName) > 0 {
				env = v.envName + "=" + os.Getenv(v.envName)
			}
			opt := strings.Join(oneArgs, ",")

			if h.MaxNameLen < len(opt) {
				h.MaxNameLen = len(opt)
			}

			switch v.pointer.Kind() {
			case reflect.Bool:
				h.Flags = append(h.Flags, showOption{Opt: opt, Usage: v.usage, Env: env})
			default:
				h.Options = append(h.Options, showOption{Opt: opt, Usage: v.usage, Env: env})
			}
		}
	}

	saveHelp(c.short)
	saveHelp(c.long)
}

func (c *Clop) printHelpMessage() {
	h := Help{}

	c.genHelpMessage(&h)

	err := h.output(os.Stdout)
	if err != nil {
		panic(err)
	}

}

func (c *Clop) parseTagAndSetOption(clop string, usage string, v reflect.Value) error {
	options := strings.Split(clop, ";")

	option := &Option{usage: usage, pointer: v}

	findName := false
	for _, opt := range options {
		opt = strings.TrimLeft(opt, " ")
		name := ""
		// TODO 检查name的长度
		switch {
		case strings.HasPrefix(opt, "--"):
			name = opt[2:]
			c.setOption(name, option, c.long)
			option.showLong = append(option.showLong, name)
			findName = true
		case strings.HasPrefix(opt, "-"):
			name = opt[1:]
			c.setOption(name, option, c.short)
			option.showShort = append(option.showShort, name)
			findName = true
		case strings.HasPrefix(opt, "def="):
			option.showDefValue = opt[4:]
		case strings.HasPrefix(opt, "greedy"):
			option.greedy = true
		case strings.HasPrefix(opt, "env="):
			findName = true
			option.envName = opt[4:]
			c.setOption(option.envName, option, c.env)
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
	return c.getOptionAndSet(a, index, numMinuses)
}

// 设置环境变量
func (c *Clop) bindEnv() error {
	for _, o := range c.env {
		if err := o.setEnv(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Clop) bindStruct() error {

	for i := 0; i < len(c.args); i++ {

		if err := c.parseOneOption(&i); err != nil {
			return err
		}

	}

	return c.bindEnv()
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

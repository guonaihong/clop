package clop

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"
)

var (
	ErrDuplicateOptions = errors.New("is already in use")
	ErrUsageEmpty       = errors.New("usage cannot be empty")
	ErrUnsupported      = errors.New("unsupported command")
	ErrNotFoundName     = errors.New("no command line options found")
	ErrOptionName       = errors.New("Illegal option name")
)

type Clop struct {
	//指向自己的root clop，如果设置了subcommand这个值是有意义的
	//非root Clop指向root，root Clop值为nil
	root         *Clop
	shortAndLong map[string]*Option  //存放长短选项
	checkEnv     map[string]struct{} //判断环境变量是否重复注册的
	checkArgs    map[string]struct{} //判断args是否重复注册
	envAndArgs   []*Option           //存放环境变量和args
	args         []string            //原始参数
	unparsedArgs []string            //没有解析的args参数

	about   string //about信息
	version string //版本信息

	exit       bool                   //测试需要用, -h --help 是否退出进程
	subcommand map[string]*Subcommand //子命令

	isSetSubcommand map[string]struct{} //用于查询哪个子命令被使用
	procName        string

	w io.Writer
}

// 使用递归定义，可以很轻检地解决subcommand嵌套的情况
type Subcommand struct {
	*Clop
	usage string
}

type Option struct {
	pointer      reflect.Value //存放需要修改的值的reflect.Value类型变量
	usage        string        //帮助信息
	showDefValue string        //显示默认值
	index        int           //表示参数优先级 TODO把值用起来
	envName      string        //环境变量
	argsName     string        //args变量
	greedy       bool          //贪婪模式 -H a b c 等于-H a -H b -H c

	showShort []string //help显示的短选项
	showLong  []string //help显示的长选项
}

func New(args []string) *Clop {
	return &Clop{
		shortAndLong:    make(map[string]*Option),
		checkEnv:        make(map[string]struct{}),
		checkArgs:       make(map[string]struct{}),
		isSetSubcommand: make(map[string]struct{}), //TODO后期优化下内存,只有root需要初始化
		args:            args,
		exit:            true,
		w:               os.Stdout,
	}
}

// 检查option 名字的合法性
func checkOptionName(name string) (byte, bool) {
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_' {
			continue
		}
		return c, false
	}
	return 0, true
}

func (c *Clop) SetExit(exit bool) *Clop {
	c.exit = exit
	return c
}

func (c *Clop) SetOutput(w io.Writer) *Clop {
	c.w = w
	return c
}

func (c *Clop) SetProcName(procName string) *Clop {
	c.procName = procName
	return c
}

func (c *Clop) IsSetSubcommand(subcommand string) bool {
	_, ok := c.isSetSubcommand[subcommand]
	return ok
}

func (c *Clop) setOption(name string, option *Option, m map[string]*Option, long bool) error {
	if c, ok := checkOptionName(name); !ok {
		return fmt.Errorf("%w:%s:unsupported characters found(%c)", ErrOptionName, name, c)
	}

	if _, ok := m[name]; ok {
		name = "-" + name
		if long {
			name = "-" + name
		}
		return fmt.Errorf("%s %w", name, ErrDuplicateOptions)
	}

	m[name] = option
	return nil
}

// 解析长选项
func (c *Clop) parseLong(arg string, index *int) error {
	var option *Option
	option, _ = c.shortAndLong[arg]
	if option == nil || len(arg) == 1 {
		return fmt.Errorf(`error: Found argument '--%s' which wasn't expected, or isn't valid in this context`,
			arg)
	}

	value := ""
	switch option.pointer.Kind() {
	//bool类型，不考虑false的情况
	case reflect.Bool:
		value = "true"
	default:
		_, isBoolSlice := option.pointer.Interface().([]bool)
		if isBoolSlice {
			return setBase("true", option.pointer)
		}

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

// 设置环境变量和参数
func (o *Option) setEnvAndArgs(c *Clop) (err error) {
	if len(o.envName) > 0 {
		if v, ok := os.LookupEnv(o.envName); ok {
			if o.pointer.Kind() == reflect.Bool {
				if v != "false" {
					v = "true"
				}
			}

			return setBase(v, o.pointer)
		}
	}

	if len(o.argsName) > 0 {
		if len(c.unparsedArgs) == 0 {
			//todo修饰下报错信息
			//return errors.New("unparsedargs == 0")
			return nil
		}

		value := c.unparsedArgs[0]
		switch o.pointer.Kind() {
		case reflect.Slice:
			for o.pointer.Kind() == reflect.Slice {
				setBase(value, o.pointer)
				c.unparsedArgs = c.unparsedArgs[1:]
				if len(c.unparsedArgs) == 0 {
					break
				}

				value = c.unparsedArgs[0]
			}
		default:
			if err := setBase(value, o.pointer); err != nil {
				return err
			}
			if len(c.unparsedArgs) > 0 {
				c.unparsedArgs = c.unparsedArgs[1:]
			}
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

		optionName := string(byte(a))
		option, _ = c.shortAndLong[optionName]
		if option == nil {
			return fmt.Errorf("error: Found argument '-%s' which wasn't expected, or isn't valid in this context", optionName)
		}

		find = true
		//value := "" //TODO
		_, isBoolSlice := option.pointer.Interface().([]bool)
		_, isBool := option.pointer.Interface().(bool)
		if !(isBoolSlice || isBool) {
			shortIndex++
		}

	getchar:
		for value := arg; ; {

			if len(value[shortIndex:]) > 0 {
				val := value[shortIndex:]
				if isBoolSlice || isBool {
					val = "true"
				}

				if err := setBase(val, option.pointer); err != nil {
					return err
				}

				if isBoolSlice || isBool { //比如-vvv这种情况
					break getchar
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

	if find {
		return nil
	}

	return fmt.Errorf(`error: Found argument '-%s' which wasn't expected, or isn't valid in this context`,
		arg)
}

func (c *Clop) getOptionAndSet(arg string, index *int, numMinuses int) error {
	// 输出帮助信息
	if arg == "h" || arg == "help" {
		c.Usage()
		return nil
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

	// shortAndLong多个key指向一个option,需要used map去重
	used := make(map[*Option]struct{}, len(c.shortAndLong))

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
				h.Flags = append(h.Flags, showOption{Opt: opt, Usage: v.usage, Env: env, Default: v.showDefValue})
			default:
				h.Options = append(h.Options, showOption{Opt: opt, Usage: v.usage, Env: env, Default: v.showDefValue})
			}
		}
	}

	saveHelp(c.shortAndLong)

	for _, v := range c.envAndArgs {
		opt := v.argsName
		if len(opt) == 0 && len(v.envName) > 0 {
			opt = v.envName
		}

		// args参数
		if len(opt) > 0 {
			opt = "<" + opt + ">"
		}
		if h.MaxNameLen < len(opt) {
			h.MaxNameLen = len(opt)
		}

		env := ""
		if len(v.envName) > 0 {
			env = v.envName + "=" + os.Getenv(v.envName)
		}
		h.Args = append(h.Args, showOption{Opt: opt, Usage: v.usage, Env: env})
	}

	// 子命令
	for opt, v := range c.subcommand {
		if h.MaxNameLen < len(opt) {
			h.MaxNameLen = len(opt)
		}
		h.Subcommand = append(h.Subcommand, showOption{Opt: opt, Usage: v.usage})
	}

	h.ProcessName = c.procName
	h.Version = c.version
	h.About = c.about
}

func (c *Clop) Usage() {
	c.printHelpMessage()
	if c.exit {
		os.Exit(0)
	}
}

func (c *Clop) printHelpMessage() {
	h := Help{}

	c.genHelpMessage(&h)

	err := h.output(c.w)
	if err != nil {
		panic(err)
	}

}

func (c *Clop) getRoot() (root *Clop) {
	root = c
	if c.root != nil {
		root = c.root
	}
	return root
}

func (c *Clop) parseSubcommandTag(clop string, usage string) (newClop *Clop, haveSubcommand bool) {
	options := strings.Split(clop, ";")
	for _, opt := range options {
		switch {
		case strings.HasPrefix(opt, "subcommand="):
			if c.subcommand == nil {
				c.subcommand = make(map[string]*Subcommand, 3)
			}

			name := opt[len("subcommand="):]
			newClop := New(nil)
			newClop.SetProcName(name)
			newClop.root = c.getRoot()
			c.subcommand[name] = &Subcommand{Clop: newClop, usage: usage}
			return newClop, true
		}
	}

	return nil, false
}

func (c *Clop) parseTagAndSetOption(clop string, usage string, def string, v reflect.Value) error {
	options := strings.Split(clop, ";")

	option := &Option{usage: usage, pointer: v, showDefValue: def}

	const (
		isShort = 1 << iota
		isLong
		isEnv
		isArgs
	)

	flags := 0
	for _, opt := range options {
		opt = strings.TrimLeft(opt, " ")
		if len(opt) == 0 {
			continue //跳过空值
		}
		name := ""
		// TODO 检查name的长度
		switch {
		//注册长选项
		case strings.HasPrefix(opt, "--"):
			name = opt[2:]
			if err := c.setOption(name, option, c.shortAndLong, true); err != nil {
				return err
			}
			option.showLong = append(option.showLong, name)
			flags |= isShort
			//注册短选项
		case strings.HasPrefix(opt, "-"):
			name = opt[1:]
			if err := c.setOption(name, option, c.shortAndLong, false); err != nil {
				return err
			}
			option.showShort = append(option.showShort, name)
			flags |= isLong
		case strings.HasPrefix(opt, "greedy"):
			option.greedy = true
		case strings.HasPrefix(opt, "env="):
			flags |= isEnv
			option.envName = opt[4:]
			if _, ok := c.checkEnv[option.envName]; ok {
				return fmt.Errorf("%s: env=%s", ErrDuplicateOptions, option.envName)
			}
			c.envAndArgs = append(c.envAndArgs, option)
			c.checkEnv[option.envName] = struct{}{}
		case strings.HasPrefix(opt, "args="):
			// args是和长,短选项互斥的功能
			if flags&isShort > 0 || flags&isLong > 0 {
				continue
			}

			flags |= isArgs
			option.argsName = opt[5:]
			if _, ok := c.checkArgs[option.argsName]; ok {
				return fmt.Errorf("%s: args=%s", ErrDuplicateOptions, option.argsName)
			}

			c.checkArgs[option.argsName] = struct{}{}
			c.envAndArgs = append(c.envAndArgs, option)

		default:
			return fmt.Errorf("%s:(%s) clop(%s)", ErrUnsupported, opt, clop)
		}

		if strings.HasPrefix(opt, "-") && len(name) == 0 {
			return fmt.Errorf("Illegal command line option:%s", opt)
		}

	}

	if flags&isShort == 0 && flags&isLong == 0 && flags&isEnv == 0 && flags&isArgs == 0 {
		return fmt.Errorf("%s:%s", ErrNotFoundName, clop)
	}

	return nil
}

func (c *Clop) registerCore(v reflect.Value, sf reflect.StructField) error {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	clop := Tag(sf.Tag).Get("clop")
	usage := Tag(sf.Tag).Get("usage")

	// 如果是subcommand
	if v.Kind() == reflect.Struct {
		if len(clop) != 0 {
			if newClop, b := c.parseSubcommandTag(clop, usage); b {
				c = newClop
			}
		}
	}

	if v.Kind() != reflect.Struct {

		def := Tag(sf.Tag).Get("default")
		if len(def) > 0 {
			if err := setDefaultValue(def, v); err != nil {
				return err
			}
		}

		if len(clop) == 0 && len(usage) == 0 {
			return nil
		}

		// 如果是存放version的字段
		if strings.HasPrefix(clop, "version=") {
			c.version = clop[8:]
			return nil
		}

		// 如果是存放about的字段
		if strings.HasPrefix(clop, "about=") {
			c.about = clop[6:]
			return nil
		}

		// clop 可以省略
		if len(usage) > 0 {
			if len(clop) == 0 {
				lowerClop := strings.ToLower(sf.Name)
				clop = "-" + string(lowerClop[0])
				if len(lowerClop) > 1 {
					clop = clop + ";--" + lowerClop
				}
			}
		}

		return c.parseTagAndSetOption(clop, usage, def, v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := typ.Field(i)

		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}

		//fmt.Printf("my.index(%d)(1.%s)-->(2.%s)\n", i, Tag(sf.Tag).Get("clop"), Tag(sf.Tag).Get("usage"))
		//fmt.Printf("stdlib.index(%d)(1.%s)-->(2.%s)\n", i, sf.Tag.Get("clop"), sf.Tag.Get("usage"))
		if err := c.registerCore(v.Field(i), sf); err != nil {
			return err
		}
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

	return c.registerCore(v, emptyField)
}

func (c *Clop) parseOneOption(index *int) error {

	arg := c.args[*index]

	if len(arg) == 0 {
		return errors.New("fail option")
	}

	if arg[0] != '-' {
		if len(c.subcommand) > 0 {
			newClop, ok := c.subcommand[arg]
			// 子命令和args都是没有-号开头，没有设置env或args就当是没有注册过的子命令，直接报错
			if !ok && len(c.envAndArgs) == 0 {
				return fmt.Errorf("Unknown subcommand:%s", arg)
			}

			c.getRoot().isSetSubcommand[arg] = struct{}{}
			newClop.args = c.args[*index+1:]
			c.args = c.args[0:0]
			return newClop.bindStruct()
		}
		c.unparsedArgs = append(c.unparsedArgs, arg)
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
func (c *Clop) bindEnvAndArgs() error {
	for _, o := range c.envAndArgs {
		if err := o.setEnvAndArgs(c); err != nil {
			return err
		}
	}

	return nil
}

// bind结构体
func (c *Clop) bindStruct() error {

	for i := 0; i < len(c.args); i++ {

		if err := c.parseOneOption(&i); err != nil {
			return err
		}

	}

	return c.bindEnvAndArgs()
}

func (c *Clop) Bind(x interface{}) (err error) {
	defer func() {
		if err != nil {
			fmt.Fprintln(c.w, err)
			fmt.Fprintln(c.w, "For more information try --help")
			if c.exit {
				os.Exit(1)
			}
		}
	}()

	if err = c.register(x); err != nil {
		return err
	}

	if err = c.bindStruct(); err != nil {
		return err
	}

	err = valid.ValidateStruct(x)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			return errors.New(e.Translate(valid.trans))
		}

	}
	return err
}

func Usage() {
	CommandLine.Usage()
}

func Bind(x interface{}) error {
	CommandLine.SetProcName(os.Args[0])
	return CommandLine.Bind(x)
}

func IsSetSubcommand(subcommand string) bool {
	return CommandLine.IsSetSubcommand(subcommand)
}

var CommandLine = New(os.Args[1:])

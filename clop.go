package clop

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var (
	ErrDuplicateOptions = errors.New("is already in use")
	//ErrUsageEmpty       = errors.New("usage cannot be empty")
	ErrUnsupported  = errors.New("unsupported clop command")
	ErrNotFoundName = errors.New("no command line options found")
	ErrOptionName   = errors.New("Illegal option name")
)

var (
	// 显示usage信息里面的[default: xxx]信息，如果为false，就不显示
	ShowUsageDefault = true
)

const (
	defautlCallbackName = "Parse"
	defaultSubMain      = "SubMain"
)

const (
	optGreedy          = "greedy"
	optOnce            = "once"
	optEnv             = "env"
	optEnvEqual        = "env="
	optSubcommand      = "subcommand"
	optSubcommandEqual = "subcommand="
	optShort           = "short"
	optLong            = "long"
	optCallback        = "callback"
	optCallbackEqual   = "callback="
	optSpace           = " "
)

/*
type SubMain interface {
	SubMain()
}
*/

type unparsedArg struct {
	arg   string
	index int
}

type Clop struct {
	//指向自己的root clop，如果设置了subcommand这个值是有意义的
	//非root Clop指向root，root Clop值为nil
	root         *Clop
	shortAndLong map[string]*Option       //存放长短选项
	checkEnv     map[string]struct{}      //判断环境变量是否重复注册的
	checkArgs    map[string]struct{}      //判断args是否重复注册
	envAndArgs   []*Option                //存放环境变量和args
	args         []string                 //原始参数
	unparsedArgs []unparsedArg            //没有解析的args参数
	allStruct    map[interface{}]struct{} //所有注册过的结构体

	about         string  //about信息
	version       string  //版本信息
	versionOption *Option //版本选项

	subMain    reflect.Value //子命令带SubMain方法, 就会自动调用
	structAddr reflect.Value
	exit       bool                   //测试需要用, 控制-h --help 是否退出进程
	subcommand map[string]*Subcommand //子命令, 保存结构体当层的所有子命令的信息

	isSetSubcommand map[string]struct{} //用于查询哪个子命令被使用, 只有root节点会设置值
	procName        string              //进程名

	currSubcommandFieldName string //当前使用的子命令结构体名, 只有root才设置该字段
	fieldName               string //记录当前子结构体字段名, root为空
	w                       io.Writer
}

// 设置版本相关信息
func (c *Clop) SetVersion(version string) *Clop {
	c.version = version
	if c.versionOption == nil {
		c.SetVersionOption("V", "version")
	}

	return c
}

// 设置版本相关信息
func (c *Clop) SetVersionOption(short, long string) *Clop {
	var opt Option

	if short != "" {
		opt.showShort = []string{short}
	}
	if long != "" {
		opt.showLong = []string{long}
	}

	c.versionOption = &opt

	return c
}

func (c *Clop) versionShort() string {
	if c.versionOption != nil && len(c.versionOption.showShort) > 0 {
		return c.versionOption.showShort[0]
	}
	return ""
}

func (c *Clop) versionLong() string {
	if c.versionOption != nil && len(c.versionOption.showLong) > 0 {
		return c.versionOption.showLong[0]
	}
	return ""
}

// 设置about相关信息
func (c *Clop) SetAbout(about string) *Clop {
	c.about = about
	return c
}

// 使用递归定义，可以很轻松地解决subcommand嵌套的情况
type Subcommand struct {
	*Clop
	usage string
}

type Option struct {
	pointer      reflect.Value //存放需要修改的值的reflect.Value类型变量
	fn           reflect.Value
	usage        string //帮助信息
	showDefValue string //显示默认值
	//表示参数优先级, 高4字节存放args顺序, 低4字节存放命令组合的顺序(ls -ltr)，这里的l的高4字节的值就是0
	index    uint64
	envName  string //环境变量
	argsName string //args变量
	greedy   bool   //贪婪模式 -H a b c 等于-H a -H b -H c
	// 如果设置once标记，命令行传递-debug -debug这种重复选项会报错
	// 对slice变量无效
	once bool //只能设置一次，如果设置once标记，命令行传了两次选项会报错

	cmdSet bool //是否通过命令行设置过值

	showShort []string //help显示的短选项
	showLong  []string //help显示的长选项
}

func (o *Option) onceResetValue() {
	if len(o.showDefValue) > 0 && !o.pointer.IsZero() && !o.cmdSet {
		resetValue(o.pointer)
	}

	o.cmdSet = true
}

func New(args []string) *Clop {
	return &Clop{
		shortAndLong:    make(map[string]*Option),
		checkEnv:        make(map[string]struct{}),
		checkArgs:       make(map[string]struct{}),
		isSetSubcommand: make(map[string]struct{}), //TODO后期优化下内存,只有root需要初始化
		allStruct:       make(map[interface{}]struct{}),
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

// 设置出错行为，默认出错会退出进程(true), 为false则不会
func (c *Clop) SetExit(exit bool) *Clop {
	c.exit = exit
	return c
}

func (c *Clop) SetOutput(w io.Writer) *Clop {
	c.w = w
	return c
}

// 设置进程名
func (c *Clop) SetProcName(procName string) *Clop {
	c.procName = procName
	return c
}

func (c *Clop) IsSetSubcommand(subcommand string) bool {
	_, ok := c.isSetSubcommand[subcommand]
	return ok
}

func (c *Clop) GetIndex(optName string) uint64 {
	// 长短选项
	o, ok := c.shortAndLong[optName]
	if ok {
		return o.index
	}

	// args参数
	if _, ok := c.checkArgs[optName]; ok {
		for _, o := range c.envAndArgs {
			if o.argsName == optName {
				return o.index
			}
		}
	}

	return 0
}

func (c *Clop) setOption(name string, option *Option, m map[string]*Option, long bool) error {
	if c, ok := checkOptionName(name); !ok {
		return fmt.Errorf("%w:%s:unsupported characters found(%c)", ErrOptionName, name, c)
	}

	if c.version != "" && (name == c.versionShort() || name == c.versionLong()) {
		name = "-" + name
		if long {
			name = "-" + name
		}

		versionOpt := ""
		if c.versionShort() != "" {
			versionOpt = "-" + c.versionShort()
		}
		if c.versionLong() != "" {
			if versionOpt != "" {
				versionOpt += ","
			}
			versionOpt += "--" + c.versionLong()
		}
		return fmt.Errorf("%s %w, duplicate definition with version option %s", name, ErrDuplicateOptions, versionOpt)
	}

	if o, ok := m[name]; ok {
		name = "-" + name
		if long {
			name = "-" + name
		}
		return fmt.Errorf("%s %w, duplicate definition with %s", name, ErrDuplicateOptions, c.showShortAndLong(o))
	}

	m[name] = option
	return nil
}

func setValueAndIndex(val string, option *Option, index int, lowIndex int) error {
	option.onceResetValue()
	option.index = uint64(index) << 31
	option.index |= uint64(lowIndex)
	if option.fn.IsValid() {
		// 如果定义callback, 就不会走默认形为
		option.fn.Call([]reflect.Value{reflect.ValueOf(val)})
		return nil
	}

	return setBase(val, option.pointer)
}

func errOnce(optionName string) error {
	return fmt.Errorf(`error: The argument '-%s' was provided more than once, but cannot be used multiple times`,
		optionName)
}

func (c *Clop) unknownOptionErrorShort(optionName string, arg string) error {
	m := fmt.Sprintf(`error: Found argument '-%s' which wasn't expected, or isn't valid in this context`,
		optionName)

	m += c.genMaybeHelpMsg(arg)
	return errors.New(m)
}

func (c *Clop) unknownOptionError(optionName string) error {
	m := fmt.Sprintf(`error: Found argument '--%s' which wasn't expected, or isn't valid in this context`,
		optionName)

	m += c.genMaybeHelpMsg(optionName)
	return errors.New(m)
}

func setBoolAndBoolSliceDefval(pointer reflect.Value, value *string) {
	kind := pointer.Kind()
	//bool类型，不考虑false的情况
	if *value == "" {
		if reflect.Bool == kind {
			*value = "true"
			return
		}

		if _, isBoolSlice := pointer.Interface().([]bool); isBoolSlice {
			*value = "true"
		}
	}
}

func (c *Clop) parseEqualValue(arg string) (value string, option *Option, err error) {
	pos := strings.Index(arg, "=")
	if pos == -1 {
		return "", nil, c.unknownOptionError(arg)
	}

	option, _ = c.shortAndLong[arg[:pos]]
	if option == nil {
		return "", nil, c.unknownOptionError(arg)
	}
	value = arg[pos+1:]
	return value, option, nil
}

func checkOnce(arg string, option *Option) error {
	if option.once && !option.pointer.IsZero() {
		return errOnce(arg)
	}
	return nil
}

func (c *Clop) isRegisterOptions(arg string) bool {
	num := 0
	if len(arg) > 0 && arg[0] == '-' {
		num++
	}

	if len(arg) > 1 && arg[1] == '-' {
		num++
	}

	// 处理value带=的情况
	end := len(arg)
	if e := strings.IndexByte(arg, '='); e != -1 {
		end = e
	}

	_, ok := c.shortAndLong[arg[num:end]]
	return ok
}

// 解析长选项
func (c *Clop) parseLong(arg string, index *int) (err error) {
	var option *Option
	value := ""
	option, _ = c.shortAndLong[arg]
	if option == nil {
		if value, option, err = c.parseEqualValue(arg); err != nil {
			return err
		}
	}

	if len(arg) == 1 {
		return c.unknownOptionError(arg)
	}

	// 设置bool 和bool slice的默认值
	setBoolAndBoolSliceDefval(option.pointer, &value)

	if len(value) > 0 {
		if err := checkOnce(arg, option); err != nil {
			return err
		}
		return setValueAndIndex(value, option, *index, 0)
	}

	// 如果是长选项
	if *index+1 >= len(c.args) {
		return nil
	}

	for {

		(*index)++
		if *index >= len(c.args) {
			return nil
		}

		value = c.args[*index]

		if c.findFallbackOpt(value, index) {
			return nil
		}

		if err := checkOnce(arg, option); err != nil {
			return err
		}

		if err := setValueAndIndex(value, option, *index, 0); err != nil {
			return err
		}

		/*
			if option.pointer.Kind() != reflect.Slice && !option.greedy {
				return nil
			}
		*/

		if !option.greedy {
			return nil
		}
	}
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

			return setValueAndIndex(v, o, 0, 0)
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
				setValueAndIndex(value.arg, o, value.index, 0)
				c.unparsedArgs = c.unparsedArgs[1:]
				if len(c.unparsedArgs) == 0 {
					break
				}

				value = c.unparsedArgs[0]
			}
		default:
			if err := setValueAndIndex(value.arg, o, value.index, 0); err != nil {
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
	// 可以解析的参数类型举例
	// -d -d是bool类型
	// -vvv 是[]bool类型
	// -d=false -d 是bool false是value
	// -ffile -f是string类型，file是value
	for shortIndex, a = range arg {
		//只支持ascii
		if a >= utf8.RuneSelf {
			return errors.New("Illegal character set")
		}

		optionName := string(byte(a))
		option, _ = c.shortAndLong[optionName]
		if option == nil {
			//没有注册过的选项直接报错
			return c.unknownOptionErrorShort(optionName, arg)
		}

		find = true
		findEqual := false //是否找到等于号
		value := arg
		_, isBoolSlice := option.pointer.Interface().([]bool)
		_, isBool := option.pointer.Interface().(bool)
		if !(isBoolSlice || isBool) {
			shortIndex++
		}

		if len(value[shortIndex:]) > 0 && len(value[shortIndex+1:]) > 0 {
			if value[shortIndex:][0] == '=' {
				findEqual = true
				shortIndex++
			}

			if value[shortIndex+1:][0] == '=' {
				findEqual = true
				shortIndex += 2
			}
		}

	getchar:
		for value := arg; ; {

			if len(value[shortIndex:]) > 0 { // 如果没有值，要取args下个参数
				val := value[shortIndex:]
				if isBoolSlice || isBool {
					val = "true"
				}

				if findEqual {
					val = string(value[shortIndex:])
				}

				if err := checkOnce(value[shortIndex:], option); err != nil {
					return err
				}

				if err := setValueAndIndex(val, option, *index, shortIndex); err != nil {
					return err
				}

				if findEqual {
					return nil
				}

				if isBoolSlice || isBool { //比如-vvv这种情况
					break getchar
				}

				/*
					非贪婪模式，解析设置slice变量，会多吃掉args参数要的变量
					if option.pointer.Kind() != reflect.Slice && !option.greedy {
						return nil
					}
				*/
				if !option.greedy {
					return nil
				}
			}

			shortIndex = 0

			if *index+1 >= len(c.args) {
				return nil
			}
			(*index)++

			value = c.args[*index]

			if c.findFallbackOpt(value, index) {
				return nil
			}

		}

	}

	if find {
		return nil
	}

	return c.unknownOptionErrorShort(arg, arg)
}

func (c *Clop) findFallbackOpt(value string, index *int) bool {

	// 如果打开贪婪模式，直到遇到-或者最后一个字符才结束
	if strings.HasPrefix(value, "-") {
		// 如果这个是命令行选项，而不是负数, 就直接回退选项
		if c.isRegisterOptions(value) {
			(*index)-- //回退这个选项
			return true
		}
	}

	return false
}

func (c *Clop) getOptionAndSet(arg string, index *int, numMinuses int) error {
	// 输出帮助信息
	if arg == "h" || arg == "help" {
		if _, ok := c.shortAndLong[arg]; !ok {
			c.Usage()
			return nil
		}
	}

	// 显示版本信息
	if c.version != "" && (arg == c.versionShort() || arg == c.versionLong()) {
		c.showVersion()
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

// ENV_NAME=
// ENV_NAME
func (o *Option) genShowEnvNameValue() (env string) {
	if len(o.envName) > 0 {
		envValue := os.Getenv(o.envName)
		env = o.envName
		if len(envValue) > 0 {
			env = env + "=" + envValue
		}
	}
	return
}

func (c *Clop) showShortAndLong(v *Option) string {
	var oneArgs []string

	for _, v := range v.showShort {
		oneArgs = append(oneArgs, "-"+v)
	}

	for _, v := range v.showLong {
		oneArgs = append(oneArgs, "--"+v)
	}
	return strings.Join(oneArgs, ",")
}

func (c *Clop) genHelpMessage(h *Help) {
	// shortAndLong多个key指向一个option,需要used map去重
	used := make(map[*Option]struct{}, len(c.shortAndLong))

	if c.shortAndLong["h"] == nil && c.shortAndLong["help"] == nil {
		c.shortAndLong["h"] = &Option{usage: "print the help information", showShort: []string{"h"}, showLong: []string{"help"}}
	}

	if c.version != "" {
		if c.versionShort() != "" && c.versionLong() != "" {
			c.shortAndLong[c.versionShort()] = &Option{usage: "print version information", showShort: []string{c.versionShort()}, showLong: []string{c.versionLong()}}
		} else if c.versionShort() != "" {
			c.shortAndLong[c.versionShort()] = &Option{usage: "print version information", showShort: []string{c.versionShort()}}
		} else {
			c.shortAndLong[c.versionLong()] = &Option{usage: "print version information", showLong: []string{c.versionLong()}}
		}
	}

	saveHelp := func(options map[string]*Option) {
		for _, v := range options {
			if _, ok := used[v]; ok {
				continue
			}

			used[v] = struct{}{}

			// 环境变量
			env := v.genShowEnvNameValue()

			opt := c.showShortAndLong(v)

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
		oldOpt := opt
		if len(opt) > 0 {
			opt = "<" + opt + ">"
		}
		if h.MaxNameLen < len(opt) {
			h.MaxNameLen = len(opt)
		}

		env := v.genShowEnvNameValue()
		if len(env) > 0 {
			h.Envs = append(h.Envs, showOption{Opt: oldOpt, Usage: v.usage, Env: env})
			continue
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
	h.ShowUsageDefault = ShowUsageDefault
}

// 显示version信息
func (c *Clop) showVersion() {
	fmt.Fprintln(c.w, c.version)
	if c.exit {
		os.Exit(0)
	}
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

func (c *Clop) parseSubcommandTag(clop string, v reflect.Value, usage string, fieldName string) (newClop *Clop, haveSubcommand bool) {
	options := strings.Split(clop, ";")
	for _, opt := range options {
		var name string
		switch {
		case strings.HasPrefix(opt, optSubcommandEqual):
			name = opt[len(optSubcommandEqual):]
		case opt == optSubcommand:
			name = strings.ToLower(fieldName)
		}
		if name != "" {
			if c.subcommand == nil {
				c.subcommand = make(map[string]*Subcommand, 3)
			}

			newClop := New(nil)
			//newClop.exit = c.exit //继承exit属性
			newClop.SetProcName(name)
			newClop.root = c.getRoot()
			c.subcommand[name] = &Subcommand{Clop: newClop, usage: usage}
			newClop.fieldName = fieldName

			newClop.subMain = v.Addr().MethodByName(defaultSubMain)
			return newClop, true
		}
	}

	return nil, false
}

func (c *Clop) parseTagAndSetOption(clop string, usage string, def string, fieldName string, v reflect.Value) (err error) {
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
		opt = strings.TrimLeft(opt, optSpace)
		if len(opt) == 0 {
			continue //跳过空值
		}
		name := ""
		// TODO 检查name的长度
		switch {
		case strings.HasPrefix(opt, optCallback):
			funcName := defautlCallbackName
			if strings.HasPrefix(opt, optCallbackEqual) {
				funcName = opt[len(optCallbackEqual):]
			}
			option.fn = c.structAddr.MethodByName(funcName)
			// 检查callback的参数长度
			if option.fn.Type().NumIn() != 1 {
				panic(fmt.Sprintf("Required function parameters->%s(val string)", funcName))
			}

		//注册长选项 --name
		case strings.HasPrefix(opt, "--"):
			name = opt[2:]
			fallthrough
		case strings.HasPrefix(opt, optLong):
			if !strings.HasPrefix(opt, "--") {
				if name, err = gnuOptionName(fieldName); err != nil {
					return err
				}
			}

			if err := c.setOption(name, option, c.shortAndLong, true); err != nil {
				return err
			}
			option.showLong = append(option.showLong, name)
			flags |= isShort
			//注册短选项
		case strings.HasPrefix(opt, "-"):
			name = opt[1:]
			fallthrough
		case strings.HasPrefix(opt, optShort):
			if !strings.HasPrefix(opt, "-") {
				if name, err = gnuOptionName(fieldName); err != nil {
					return err
				}
				name = string(name[0])
			}

			if err := c.setOption(name, option, c.shortAndLong, false); err != nil {
				return err
			}
			option.showShort = append(option.showShort, name)
			flags |= isLong
		case strings.HasPrefix(opt, optGreedy):
			option.greedy = true
		case strings.HasPrefix(opt, optOnce):
			option.once = true
		case opt == optEnv:
			if name, err = envOptionName(fieldName); err != nil {
				return err
			}
			fallthrough
		case strings.HasPrefix(opt, optEnvEqual):
			flags |= isEnv
			if strings.HasPrefix(opt, optEnvEqual) {
				name = opt[4:]
			}

			option.envName = name
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
			return fmt.Errorf(`%s:(%s) clop:"%s", Maybe you need to clop:"short;long"`, ErrUnsupported, opt, clop)
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
			if newClop, b := c.parseSubcommandTag(clop, v, usage, sf.Name); b {
				c = newClop
			}
		}
	}

	if v.Kind() != reflect.Struct {
		def := Tag(sf.Tag).Get("default")
		def = strings.TrimSpace(def)
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

		return c.parseTagAndSetOption(clop, usage, def, sf.Name, v)
	}

	typ := v.Type()
	c.structAddr = v.Addr()
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
	if x == nil {
		return ErrUnsupportedType
	}

	v := reflect.ValueOf(x)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("%s:got(%T)", ErrNotPointerType, v.Type())
	}

	// 如果v不是指针 v.IsNil()函数调用会崩溃，所以指针要放到前面判断
	if v.IsNil() {
		return ErrUnsupportedType
	}

	c.allStruct[x] = struct{}{}
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
			if c.root == nil {
				c.currSubcommandFieldName = newClop.fieldName
			}

			newClop.args = c.args[*index+1:]
			c.args = c.args[0:0]
			err := newClop.bindStruct()
			if err != nil {
				return err
			}
			if newClop.subMain.IsValid() {
				newClop.subMain.Call([]reflect.Value{})
			}
		}
		c.unparsedArgs = append(c.unparsedArgs, unparsedArg{arg: arg, index: *index})
		return nil
	}

	// arg 必须是减号开头的字符串
	numMinuses := 1

	if arg == "-" {
		c.unparsedArgs = append(c.unparsedArgs, unparsedArg{arg: arg, index: *index})
		return nil
	}

	if arg[1] == '-' {
		numMinuses++
	}

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

	if len(c.currSubcommandFieldName) > 0 {
		v := reflect.ValueOf(x)
		v = v.Elem() // x只能是指针，已经在c.register判断过了
		v = v.FieldByName(c.currSubcommandFieldName)
		// 只有设置过的子命令才需要数据校验
		// 这里把root结构体删除掉
		delete(c.allStruct, x)
		x = v.Addr().Interface()
	}

	c.allStruct[x] = struct{}{}

	for x := range c.allStruct {
		err = valid.ValidateStruct(x)
		if err != nil {
			errs := err.(validator.ValidationErrors)

			for _, e := range errs {
				// can translate each error one at a time.
				return errors.New(e.Translate(valid.trans))
			}

		}
	}
	return err
}

// MustBind 和Bind函数类似， 出错直接panic
func (c *Clop) MustBind(x interface{}) {
	if err := c.Bind(x); err != nil {
		panic(err.Error())
	}
}

// 只注册结构体信息, 不解析
func (c *Clop) Register(x interface{}) error {
	return c.register(x)
}

// 打印帮助信息
func Usage() {
	CommandLine.Usage()
}

func MustRegister(x interface{}) {
	err := CommandLine.Register(x)
	if err != nil {
		panic(err.Error())
	}
}

// Bind接口, 包含以下功能
// 结构体字段注册
// 命令行解析
func Bind(x interface{}) error {
	CommandLine.SetProcName(os.Args[0])
	return CommandLine.Bind(x)
}

// 设置版本号
func SetVersion(version string) {
	CommandLine.SetVersion(version)
}

// 设置版本号选项，覆盖默认的V和Version
func SetVersionOption(short, long string) {
	CommandLine.SetVersionOption(short, long)
}

func SetAbout(about string) {
	CommandLine.SetAbout(about)
}

// Bind必须成功的版本
func MustBind(x interface{}) {
	CommandLine.SetProcName(os.Args[0])
	CommandLine.MustBind(x)
}

func IsSetSubcommand(subcommand string) bool {
	return CommandLine.IsSetSubcommand(subcommand)
}

func GetIndex(optName string) uint64 {
	return CommandLine.GetIndex(optName)
}

var CommandLine = New(os.Args[1:])

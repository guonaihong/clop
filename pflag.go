package clop

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
)

// flag库的解析函数名, 白名称
type argsNumAndType struct {
	size     int
	typeName string
}

var funcName = map[string]argsNumAndType{
	//"Func":true, TODO
	"Bool":        {3, "bool"},
	"BoolVar":     {4, "bool"},
	"Duration":    {3, "time.Duration"},
	"DurationVar": {4, "time.Duration"},
	"Float64":     {3, "float64"},
	"Float64Var":  {4, "float64"},
	"Int":         {3, "int"},
	"IntVar":      {4, "int"},
	"Int64":       {3, "int64"},
	"Int64Var":    {4, "int64"},
	"String":      {3, "string"},
	"StringVar":   {4, "string"},
	"Uint":        {3, "uint"},
	"UintVar":     {4, "uint"},
	"Uint64":      {3, "uint64"},
	"Uint64Var":   {4, "uint64"},
}

// 解析flag
type ParseFlag struct {
	astFile     *ast.File
	fileName    string
	src         string
	funcAndArgs map[string]funcAndArgs
	allOutBuf   bytes.Buffer
}

// 构造函数
func NewParseFlag() *ParseFlag {
	return &ParseFlag{}
}

// 参数
type funcAndArgs struct {
	args      []flagOpt
	haveParse bool
}

// 保存从ast里面提取出来的元数据
type flagOpt struct {
	varName  string
	optName  string
	defVal   string
	usage    string
	typeName string
}

// 可以判断是你要的函数, 比如flag.String
func isFunc(expr ast.Expr, pkg, fn string) bool {
	f, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(f.X, pkg) && isIdent(f.Sel, fn)
}

func getPtrArgName(arg ast.Expr) string {
	a, ok := arg.(*ast.UnaryExpr)
	if !ok {
		return ""
	}
	return getIdentName(a.X)
}

// 获取函数名
func getArgName(arg ast.Expr) string {
	a, ok := arg.(*ast.BasicLit)
	if !ok {
		return ""
	}
	return a.Value
}

// 提取函数名和形参
func (p *ParseFlag) takeFuncNameAndArgs(expr ast.Expr, args []ast.Expr) {
	f, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return
	}

	obj := getIdentName(f.X)
	fn := getIdentName(f.Sel)

	argsNumType, ok := funcName[fn]
	if !ok {
		return
	}

	if _, ok := p.funcAndArgs[obj]; !ok {
		p.funcAndArgs[obj] = funcAndArgs{}
	}

	if argsNumType.size != len(args) {
		return
	}

	var opt flagOpt
	if argsNumType.size == 3 {
		opt.varName = obj
		opt.optName = getArgName(args[0])
		opt.defVal = getArgName(args[1])
		opt.usage = getArgName(args[2])

	} else {
		opt.varName = getPtrArgName(args[0])
		opt.optName = getArgName(args[1])
		opt.defVal = getArgName(args[2])
		opt.usage = getArgName(args[3])
	}

	opt.typeName = argsNumType.typeName
	oldVal := p.funcAndArgs[fn]
	oldVal.args = append(oldVal.args, opt)
	p.funcAndArgs[fn] = oldVal

}

func isIdent(expr ast.Expr, name string) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == name
}

func getIdentName(expr ast.Expr) string {
	ident, ok := expr.(*ast.Ident)
	if ok {
		return ident.Name
	}

	return ""
}

// 解析flag.NewFlagSet代码
func (p *ParseFlag) parserFlagNewFlagSet(stmt *ast.AssignStmt) {
	if (stmt.Tok == token.ASSIGN || stmt.Tok == token.DEFINE) && len(stmt.Rhs) > 0 {
		if isFunc(stmt.Rhs[0], "flag", "NewFlagSet") && len(stmt.Lhs) > 0 {
			name := getIdentName(stmt.Lhs[0])
			if len(name) == 0 {
				return
			}

			p.funcAndArgs[name] = funcAndArgs{}

		}
	}
}

func parserFlagParser() {
}

// 解析函数调用的代码
func (p *ParseFlag) findFuncCalls(node ast.Node) bool {
	stmt, ok := node.(*ast.AssignStmt)
	if ok {
		p.parserFlagNewFlagSet(stmt)
		return true
	}

	call, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}

	p.takeFuncNameAndArgs(call.Fun, call.Args)

	return true
}

// 获取函数和形参
func (p *ParseFlag) getFuncCallsToken() (err error) {

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, p.fileName, p.src, 0)
	if err != nil {
		return err
	}

	p.funcAndArgs["flag"] = funcAndArgs{}
	p.astFile = f

	ast.Inspect(p.astFile, p.findFuncCalls)
	return nil
}

func (p *ParseFlag) Parse() {
	p.getFuncCallsToken()
}

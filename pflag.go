package clop

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// 解析flag
type ParseFlag struct {
	astFile  *ast.File
	fileName string
	src      string
}

// 构造函数
func NewParseFlag() *ParseFlag {
	return &ParseFlag{}
}

// 保存从ast里面提取出来的元数据
type flagOpt struct {
	varName string
	optName string
	defVal  string
	usage   string
}

func isFunc(expr ast.Expr, pkg, fn string) bool {
	f, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(f.X, pkg) && isIdent(f.Sel, fn)
}

func isIdent(expr ast.Expr, name string) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == name
}

func (p *ParseFlag) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), p.astFile)
}

func (p *ParseFlag) findFuncCalls(node ast.Node) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return true
	}

	_ = call
	//isFunc(call, )

	return true
}

func (p *ParseFlag) funcCallsToken() (err error) {

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, p.fileName, p.src, 0)
	if err != nil {
		return err
	}

	p.astFile = f

	p.walk(p.findFuncCalls)
	return nil
}

func (p *ParseFlag) Parse() {
	p.funcCallsToken()
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}

	return nil
}

package clop

import (
	"go/ast"
	"go/parser"
	"go/token"
)

var funcName = map[string]bool{
	//"Func":true, TODO
	"Bool":        true,
	"BoolVar":     true,
	"Duration":    true,
	"DurationVar": true,
	"Float64":     true,
	"Float64Var":  true,
	"Int":         true,
	"IntVar":      true,
	"Int64":       true,
	"Int64Var":    true,
	"String":      true,
	"StringVar":   true,
	"Uint":        true,
	"UintVar":     true,
	"Uint64":      true,
	"Uint64Var":   true,
}

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

func parserFlagNewFlagSet(stmt *ast.AssignStmt) {
	if (stmt.Tok == token.ASSIGN || stmt.Tok == token.DEFINE) && len(stmt.Rhs) > 0 {
		isFunc(stmt.Rhs[0], "flag", "NewFlagSet")
	}
}

func (p *ParseFlag) findFuncCalls(node ast.Node) bool {
	stmt, ok := node.(*ast.AssignStmt)
	if ok {
		parserFlagNewFlagSet(stmt)
		return true
	}

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

	ast.Inspect(p.astFile, p.findFuncCalls)
	return nil
}

func (p *ParseFlag) Parse() {
	p.funcCallsToken()
}

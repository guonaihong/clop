package clop

import (
	"go/ast"
)

// 解析flag
type ParseFlag struct {
	astFile *ast.File
}

// 保存从ast里面提取出来的元数据
type flagOpt struct {
	varName string
	optName string
	defVal  string
	usage   string
}

func (p *ParseFlag) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), p.astFile)
}

func (p *ParseFlag) funcCalls() {
}

func (p *ParseFlag) FromFile() {
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	return w
}

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"reflect"
)

var m map[reflect.Type]func(ast.Node)

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {

	if node != nil {
		// fmt.Println(reflect.TypeOf(node))
		if f := m[reflect.TypeOf(node)]; f != nil {
			f(node)
		}
		// fmt.Println(fset.Position(node.Pos()))
		// fmt.Println(fset.Position(node.End()))
	}

	return v
}

var fset *token.FileSet

func main() {
	fset = token.NewFileSet() // positions are relative to fset
	m = make(map[reflect.Type]func(ast.Node))

	m[reflect.TypeOf((*ast.FuncDecl)(nil))] = func(node ast.Node) {
		f := node.(*ast.FuncDecl)
		fmt.Println(f.Body)
		f.Body = &ast.BlockStmt{token.NoPos, make([]ast.Stmt, 0), token.NoPos}
	}

	f, err := parser.ParseDir(fset, "./test_fixtures", nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for pkgName, pkg := range f {
		fmt.Println(pkgName)
		fmt.Println(pkg.Files)
		for _, file := range pkg.Files {
			ast.Walk(&Visitor{}, file)
			printer.Fprint(os.Stdout, fset, file)
		}
	}
}

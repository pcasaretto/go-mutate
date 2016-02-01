package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
)

var m map[reflect.Type]func(ast.Node)

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) (w ast.Visitor) {

	if node != nil {
		fmt.Println(fset.Position(node.Pos()))
		fmt.Println(fset.Position(node.End()))
		fmt.Println(reflect.TypeOf(node))
		if f := m[reflect.TypeOf(node)]; f != nil {
			f(node)
		}
	}

	return v
}

var fset *token.FileSet

func main() {
	fset = token.NewFileSet() // positions are relative to fset
	m = make(map[reflect.Type]func(ast.Node))

	// m[reflect.TypeOf((*ast.FuncDecl)(nil))] = func(node ast.Node) {
	// 	f := node.(*ast.FuncDecl)
	// 	fmt.Println(f.Body)
	// 	f.Body = &ast.BlockStmt{token.NoPos, make([]ast.Stmt, 0), token.NoPos}
	// }

	m[reflect.TypeOf((*ast.BasicLit)(nil))] = func(node ast.Node) {
		f := node.(*ast.BasicLit)
		fmt.Println(f.Value)
		fmt.Println(f.Kind)
		f.Value = "0"
	}

	dir := "./test_fixtures"

	f, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for pkgName, pkg := range f {
		fmt.Println(pkgName)
		fmt.Println(pkg.Files)
		for filename, file := range pkg.Files {
			if isTestFile(filename) {
				continue
			}
			ast.Walk(&Visitor{}, file)
			var w io.Writer
			f, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				w = os.Stdout
			} else {
				w = f
				defer f.Close()
			}
			printer.Fprint(w, fset, file)
		}
	}
}

// isTestFile reports whether the source file is a set of tests and should therefore
// be excluded from coverage analysis.
func isTestFile(file string) bool {
	// We don't cover tests, only the code they test.
	return strings.HasSuffix(file, "_test.go")
}

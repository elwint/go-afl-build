package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func findFuzzFunc() (string, error) {
	// Parse the Go package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, `.`, nil, 0)
	if err != nil {
		return ``, err
	}

	// Find the Fuzz function in the package
	for _, pkg := range pkgs {
		for fname, file := range pkg.Files {
			if strings.HasSuffix(fname, `_test.go`) {
				continue
			}
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == *funcName {
					return pkg.Name, nil
				}
			}
		}
	}
	return ``, fmt.Errorf("fuzz function %s not found", *funcName)
}

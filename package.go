package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

func findPackageFunc() (string, error) {
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

func loadPackage(pkgName string) (*packages.Package, error) {
	if pkgName == `main` {
		return nil, nil
	}

	pkgs, err := packages.Load(nil)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 || pkgs[0].Name != pkgName {
		return nil, fmt.Errorf(`could not load package %s`, pkgName)
	}
	return pkgs[0], nil
}

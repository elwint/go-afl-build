package main

import (
	"os"
	"text/template"

	_ "embed"

	"golang.org/x/tools/go/packages"
)

//go:embed tmpl/main.go.tmpl
var tmplMainGo string

//go:embed tmpl/main.c.tmpl
var tmplMainC string

type mainGoData struct {
	PkgImport *packages.Package
	PkgName   string
	FuncName  string
}

func createLibFile() (string, string, func()) {
	libFile := createEmptyFile(`afl.*.a`)
	libHeader := libFile[:len(libFile)-1] + `h`

	return libFile, libHeader, func() {
		os.Remove(libFile)
		os.Remove(libHeader)
	}
}

func createTemplate(tmplStr, pattern string, data any) (string, func()) {
	tmpl, err := template.New(``).Parse(tmplStr)
	panicOnError(err)

	tmpFile, err := os.CreateTemp(`.`, pattern)
	panicOnError(err)
	defer tmpFile.Close()

	cleanup := func() { os.Remove(tmpFile.Name()) }

	err = tmpl.Execute(tmpFile, data)
	if err != nil {
		cleanup()
		panic(err)
	}

	return tmpFile.Name(), cleanup
}

func createEmptyFile(pattern string) string {
	tmpFile, err := os.CreateTemp(`.`, pattern)
	if err == nil {
		err = tmpFile.Close()
	}
	panicOnError(err)

	return tmpFile.Name()
}

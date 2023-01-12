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

func createMainGo(data mainGoData) (string, func()) {
	tmpl, err := template.New(``).Parse(tmplMainGo)
	panicOnError(err)

	mainGo, cleanup := createTempFile(`main.*.go`)

	err = tmpl.Execute(mainGo, data)
	if err != nil {
		cleanup()
		panicOnError(err)
	}

	return mainGo.Name(), cleanup
}

func createLibFile() (string, func()) {
	libFileName := createEmptyFile(`afl.*.a`)

	return libFileName, func() {
		os.Remove(libFileName)
		os.Remove(libFileName[:len(libFileName)-1] + `h`)
	}
}

func createMainC() (string, func()) {
	mainC, cleanup := createTempFile(`main.*.c`)

	_, err := mainC.WriteString(tmplMainC)
	if err != nil {
		cleanup()
		panicOnError(err)
	}

	return mainC.Name(), cleanup
}

func createTempFile(pattern string) (*os.File, func()) {
	tmpFile, err := os.CreateTemp(`.`, pattern)
	if err != nil {
		panic(err)
	}
	return tmpFile, func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}
}

func createEmptyFile(pattern string) string {
	tmpFile, err := os.CreateTemp(`.`, pattern)
	if err == nil {
		err = tmpFile.Close()
	}
	if err != nil {
		panic(err)
	}

	return tmpFile.Name()
}

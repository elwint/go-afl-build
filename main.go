package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	_ "embed"

	"golang.org/x/tools/go/packages"
)

//go:embed tmpl/main.go.tmpl
var tmplMainGo string

//go:embed tmpl/main.c.tmpl
var tmplMainC string

var funcName = flag.String("func", "Fuzz", "name of the Fuzz function")

func main() {
	flag.Parse()
	if *funcName == "" {
		fmt.Println("Usage: go-afl-build -func FuncName")
		flag.PrintDefaults()
		os.Exit(1)
	}

	pkgName, err := findFuzzFunc()
	if err != nil {
		fmt.Printf("Error parsing package: %s\n", err)
		os.Exit(1)
	}

	var pkgImport *packages.Package
	if pkgName != `main` {
		pkgs, err := packages.Load(nil)
		if err != nil {
			panic(err)
		}
		if len(pkgs) == 0 || pkgs[0].Name != pkgName {
			panic(`could not load package ` + pkgName)
		}
		pkgImport = pkgs[0]
	}

	aflCC := os.Getenv(`AFL_CC`)
	if aflCC == `` {
		aflCC = `afl-gcc-fast`
	}
	gccGo := os.Getenv(`GCCGO`)
	if gccGo == `` {
		gccGo = `gccgo`
	}

	// Switch environment variables so that afl uses gccgo
	err = os.Setenv(`AFL_CC`, gccGo)
	if err != nil {
		panic(err)
	}
	err = os.Setenv(`GCCGO`, aflCC)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(``).Parse(tmplMainGo)
	if err != nil {
		panic(err)
	}

	mainGo, cleanupMainGo := createTempFile(`main.*.go`)
	defer cleanupMainGo()

	err = tmpl.Execute(mainGo, struct {
		PkgImport *packages.Package
		PkgName   string
		FuncName  string
	}{
		PkgImport: pkgImport,
		PkgName:   pkgName,
		FuncName:  *funcName,
	})
	if err != nil {
		panic(err)
	}

	libFileName := createEmptyFile(`afl.*.a`)
	defer os.Remove(libFileName)
	defer os.Remove(libFileName[:len(libFileName)-1] + `h`)

	args := []string{`build`, `-x`, `-compiler`, `gccgo`, `-buildmode`, `c-archive`, `-o`, libFileName}
	if pkgName != `main` {
		args = append(args, mainGo.Name())
	}

	out, err := exec.Command(`go`, args...).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	mainC, cleanupMainC := createTempFile(`main.*.c`)
	defer cleanupMainC()

	_, err = mainC.WriteString(tmplMainC)
	if err != nil {
		panic(err)
	}

	out, err = exec.Command(aflCC, `-o`, `afl`, mainC.Name(), libFileName).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
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

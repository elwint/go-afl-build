package main

import (
	"flag"
	"fmt"
	"os"
)

var funcName = flag.String("func", "Fuzz", "name of the Fuzz function")

func main() {
	flag.Parse()
	checkFlags()

	// Switch environment variables so that afl uses gccgo
	aflCC := setAflEnvVars()

	// Find the package and the Fuzz function name
	pkgName, err := findPackageFunc()
	handlePkgError(err)

	// Load the package
	pkgImport, err := loadPackage(pkgName)
	handlePkgError(err)

	// Create a temp main go file
	mainGo, cleanupMainGo := createTemplate(tmplMainGo, `main.*.go`, mainGoData{
		PkgImport: pkgImport,
		PkgName:   pkgName,
		FuncName:  *funcName,
	})
	defer cleanupMainGo()

	// Create a temp library file
	libFile, libHeader, cleanupLibFile := createLibFile()
	defer cleanupLibFile()

	// Build the library file using gccgo
	buildLibFile(pkgName, mainGo, libFile)

	// Create a temp main c file
	mainC, cleanupMainC := createTemplate(tmplMainC, `main.*.c`, libHeader)
	defer cleanupMainC()

	// Compile with AFL++ compiler
	buildAFL(aflCC, mainC, libFile)
}

func checkFlags() {
	if *funcName == "" {
		fmt.Println("Usage: go-afl-build -func FuncName")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setAflEnvVars() string {
	aflCC := os.Getenv(`AFL_CC`)
	if aflCC == `` {
		aflCC = `afl-gcc-fast`
	}
	gccGo := os.Getenv(`GCCGO`)
	if gccGo == `` {
		gccGo = `gccgo`
	}

	err := os.Setenv(`AFL_CC`, gccGo)
	panicOnError(err)
	err = os.Setenv(`GCCGO`, aflCC)
	panicOnError(err)

	return aflCC
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	_ "embed"
)

//go:embed tmpl/main.go.tmpl
var tmplMainGo string

//go:embed tmpl/main.c.tmpl
var tmplMainC string

func main() {
	aflCC := os.Getenv(`AFL_CC`)
	if aflCC == `` {
		aflCC = `afl-gcc-fast`
	}
	gccGo := os.Getenv(`GCCGO`)
	if gccGo == `` {
		gccGo = `gccgo`
	}

	// Switch environment variables so that afl uses gccgo
	err := os.Setenv(`AFL_CC`, gccGo)
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

	err = tmpl.Execute(mainGo, nil)
	if err != nil {
		panic(err)
	}

	libFileName := createEmptyFile(`afl.*.a`)
	defer os.Remove(libFileName)
	defer os.Remove(libFileName[:len(libFileName)-1] + `h`)

	out, err := exec.Command(`go`, `build`, `-x`, `-compiler`, `gccgo`, `-buildmode`, `c-archive`, `-o`, libFileName).CombinedOutput()
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

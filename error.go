package main

import (
	"fmt"
	"os"
)

func handlePkgError(err error) {
	if err != nil {
		fmt.Printf("Error parsing package: %s\n", err)
		os.Exit(1)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleCmdError(err error, out []byte) {
	if err != nil {
		os.Stderr.Write(out)
		panic(err)
	}
}

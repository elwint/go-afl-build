package main

import "os/exec"

func buildLibFile(pkgName, mainGo, libFile string) {
	args := []string{`build`, `-x`, `-compiler`, `gccgo`, `-buildmode`, `c-archive`, `-o`, libFile}
	if pkgName != `main` {
		args = append(args, mainGo)
	}
	out, err := exec.Command(`go`, args...).CombinedOutput()
	handleCmdError(err, out)
}

func buildAFL(aflCC, mainC, libFile string) {
	out, err := exec.Command(aflCC, `-o`, `afl`, mainC, libFile).CombinedOutput()
	handleCmdError(err, out)
}

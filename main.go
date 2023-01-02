package main

import (
	"fmt"
	"os"
	"os/exec"
)

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

	err = os.WriteFile(`go-afl-build.go`, []byte(`
package main

import (
	"C"
	"unsafe"
)

//export __FuzzAFL
func __FuzzAFL(s *C.uchar, size C.int) {
	data := (*[1 << 30]byte)(unsafe.Pointer(s))[:size:size]
	Fuzz(data)
}

func main() {}`), 0o644)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = os.Remove(`go-afl-build.go`)
	}()

	out, err := exec.Command(`go`, `build`, `-x`, `-compiler`, `gccgo`, `-buildmode`, `c-archive`, `-o`, `go-afl-build.a`).Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	defer func() {
		_ = os.Remove(`go-afl-build.a`)
		_ = os.Remove(`go-afl-build.h`)
	}()

	// https://github.com/AFLplusplus/AFLplusplus/blob/stable/instrumentation/README.persistent_mode.md
	err = os.WriteFile(`go-afl-build.c`, []byte(`
__AFL_FUZZ_INIT();
int main() {
#ifdef __AFL_HAVE_MANUAL_CONTROL
  __AFL_INIT();
#endif
  unsigned char *buf = __AFL_FUZZ_TESTCASE_BUF;
  while (__AFL_LOOP(1000)) {
    int len = __AFL_FUZZ_TESTCASE_LEN;
    __FuzzAFL(buf, len);
  }
}`), 0o644)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = os.Remove(`go-afl-build.c`)
	}()

	out, err = exec.Command(aflCC, `-o`, `afl`, `go-afl-build.c`, `go-afl-build.a`).Output()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

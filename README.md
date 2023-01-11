# go-afl-build

Wrapper to instrument Go binaries for AFL++ using gccgo and afl-gcc-fast with persistent mode

**WARNING: Highly experimental!**

example.go:

```go
package main

func Fuzz(data []byte) {
	// Call function to fuzz
}
```

```
go-afl-build -func Fuzz
afl-fuzz -i corpus -t 1000 -o output ./afl
```

**Note: the first test case will always freeze for some reason, use a timeout (e.g. `afl-fuzz -t 1000`) to skip the first test case.**

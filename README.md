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

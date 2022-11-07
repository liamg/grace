package main

import (
	"fmt"
	"syscall"
)

func main() {
	var buf syscall.Utsname
	_ = syscall.Uname(&buf)
	printUtsField(buf.Release)
}

func printUtsField(f [65]int8) {
	var str []byte
	for i := 0; i < len(f); i++ {
		str = append(str, byte(f[i]))
	}
	fmt.Println(string(str))
}

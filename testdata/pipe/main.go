package main

import (
	"fmt"
	"syscall"
)

func main() {
	fds := make([]int, 2)
	_ = syscall.Pipe2(fds, 0)
	fmt.Println(fds[0], fds[1])
}

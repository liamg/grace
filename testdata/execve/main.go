package main

import (
	"os"
	"syscall"
)

func main() {
	_ = syscall.Exec("/usr/bin/ls", []string{"ls", "/tmp"}, os.Environ())
}

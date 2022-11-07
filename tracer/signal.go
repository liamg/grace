package tracer

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type SigInfo struct {
	Signo  int32
	Errno  int32
	Code   int32
	TrapNo int32
	Pid    int32
	Uid    int32
}

func getSignalInfo(pid int) (*SigInfo, error) {
	var info SigInfo
	_, _, e1 := syscall.Syscall6(syscall.SYS_PTRACE, uintptr(unix.PTRACE_GETSIGINFO), uintptr(pid), 0, uintptr(unsafe.Pointer(&info)), 0, 0)
	if e1 != 0 {
		return nil, fmt.Errorf("ptrace get signal info failed: %v", e1)
	}
	return &info, nil
}

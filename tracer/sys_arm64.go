//go:build arm64

package tracer

import (
	"syscall"
)

const bitSize = 64

// useful info: https://chromium.googlesource.com/chromiumos/docs/+/master/constants/syscalls.md

func parseSyscall(regs *syscall.PtraceRegs) *Syscall {
	return &Syscall{
		number: int(regs.Regs[8]),
		rawArgs: [6]uintptr{
			uintptr(regs.Regs[0]),
			uintptr(regs.Regs[1]),
			uintptr(regs.Regs[2]),
			uintptr(regs.Regs[3]),
			uintptr(regs.Regs[4]),
			uintptr(regs.Regs[5]),
		},
		rawRet: uintptr(regs.Regs[0]),
	}
}

// TODO: add syscall table for arm64
var sysMap = map[int]SyscallMetadata{}

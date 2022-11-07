//go:build amd64

package tracer

import (
	"bytes"
	"syscall"

	"github.com/liamg/grace/tracer/annotation"

	"golang.org/x/sys/unix"
)

const bitSize = 64

// useful info: https://chromium.googlesource.com/chromiumos/docs/+/master/constants/syscalls.md

func parseSyscall(regs *syscall.PtraceRegs) *Syscall {
	return &Syscall{
		number: int(regs.Orig_rax),
		rawArgs: [6]uintptr{
			uintptr(regs.Rdi),
			uintptr(regs.Rsi),
			uintptr(regs.Rdx),
			uintptr(regs.R10),
			uintptr(regs.R8),
			uintptr(regs.R9),
		},
		rawRet: uintptr(regs.Rax),
	}
}

var (
	sysMap = map[int]SyscallMetadata{
		unix.SYS_READ: {
			Name: "read",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_WRITE: {
			Name: "write",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_OPEN: {
			Name: "open",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_CLOSE: {
			Name: "close",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_STAT: {
			Name: "stat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:        "stat",
					Type:        argTypeStat,
					Destination: true,
				},
			},
		},
		unix.SYS_FSTAT: {
			Name: "fstat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "stat",
					Type:        argTypeStat,
					Destination: true,
				},
			},
		},
		unix.SYS_LSTAT: {
			Name: "lstat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:        "stat",
					Type:        argTypeStat,
					Destination: true,
				},
			},
		},
		unix.SYS_POLL: {
			Name: "poll",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "ufds",
					Type:        argTypePollFdArray,
					Destination: true,
				},
				{
					Name: "nfds",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "timeout",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LSEEK: {
			Name: "lseek",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name:      "whence",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotation.AnnotateWhence,
				},
			},
		},
		unix.SYS_MMAP: {
			Name: "mmap",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeAddress,
			},
			Args: []ArgMetadata{
				{
					Name:      "addr",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "prot",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateProt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateMMapFlags,
				},
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "off",
					Type: ArgTypeUnsignedLong,
				},
			},
		},
		unix.SYS_MPROTECT: {
			Name: "mprotect",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "start",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name:      "prot",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateProt,
				},
			},
		},
		unix.SYS_MUNMAP: {
			Name: "munmap",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "start",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_BRK: {
			Name: "brk",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeAddress,
			},
			Args: []ArgMetadata{
				{
					Name:      "brk",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
			},
		},
		unix.SYS_RT_SIGACTION: {
			Name: "rt_sigaction",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "signum",
					Type: ArgTypeInt,
				},
				{
					Name: "act",
					Type: argTypeSigAction,
				},
				{
					Name:        "oldact",
					Type:        argTypeSigAction,
					Destination: true, // TODO: is this correct?
				},
			},
		},
		unix.SYS_RT_SIGPROCMASK: {
			Name: "rt_sigprocmask",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "how",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSigProcMaskFlags,
				},
				{
					Name:        "set",
					Type:        ArgTypeAddress,
					Destination: true,
				},
				{
					Name:        "oldset",
					Type:        ArgTypeAddress,
					Destination: true,
				},
			},
		},
		unix.SYS_RT_SIGRETURN: {
			Name: "rt_sigreturn",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "__unused",
					Type: ArgTypeLong,
				},
			},
		},
		unix.SYS_IOCTL: {
			Name: "ioctl",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "request",
					Type: ArgTypeInt,
				},
				{
					Name: "argp",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_PREAD64: {
			Name: "pread64",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{

				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceNext,
					Destination: true,
				},
				{
					Name: "count",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PWRITE64: {
			Name: "pwrite64",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{

				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "count",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_READV: {
			Name: "readv",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_WRITEV: {
			Name: "writev",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_ACCESS: {
			Name: "access",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_PIPE: {
			Name: "pipe",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "pipefd",
					Type:        argTypeIntArray,
					LenSource:   LenSourceFixed,
					FixedCount:  2,
					Destination: true,
				},
			},
		},
		unix.SYS_SELECT: {
			Name: "select",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "nfds",
					Type: ArgTypeInt,
				},
				{
					Name: "readfds",
					Type: argTypeFdSet,
				},
				{
					Name: "writefds",
					Type: argTypeFdSet,
				},
				{
					Name: "exceptfds",
					Type: argTypeFdSet,
				},
				{
					Name: "timeout",
					Type: argTypeTimeval,
				},
			},
		},
		unix.SYS_SCHED_YIELD: {
			Name: "sched_yield",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
		},
		unix.SYS_MREMAP: {
			Name: "mremap",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "old_address",
					Type: ArgTypeAddress,
				},
				{
					Name: "old_size",
					Type: ArgTypeInt,
				},
				{
					Name: "new_size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMRemapFlags,
				},
			},
		},
		unix.SYS_MSYNC: {
			Name: "msync",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "length",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMSyncFlags,
				},
			},
		},
		unix.SYS_MINCORE: {
			Name: "mincore",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "length",
					Type: ArgTypeInt,
				},
				{
					Name: "vec",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_MADVISE: {
			Name: "madvise",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "length",
					Type: ArgTypeInt,
				},
				{
					Name:      "advice",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMAdviseAdvice,
				},
			},
		},
		unix.SYS_SHMGET: {
			Name: "shmget",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeAddress,
			},
			Args: []ArgMetadata{
				{
					Name: "key",
					Type: ArgTypeInt,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "shmflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSHMGetFlags,
				},
			},
		},
		unix.SYS_SHMAT: {
			Name: "shmat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeAddress,
			},
			Args: []ArgMetadata{
				{
					Name: "shmid",
					Type: ArgTypeInt,
				},
				{
					Name: "shmaddr",
					Type: ArgTypeAddress,
				},
				{
					Name:      "shmflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSHMAtFlags,
				},
			},
		},
		unix.SYS_SHMCTL: {
			Name: "shmctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "shmid",
					Type: ArgTypeInt,
				},
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSHMCTLCommand,
				},
				{
					Name:        "buf",
					Type:        argTypeSHMIDDS,
					Destination: true,
				},
			},
		},
		unix.SYS_DUP: {
			Name: "dup",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldfd",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_DUP2: {
			Name: "dup2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldfd",
					Type: ArgTypeInt,
				},
				{
					Name: "newfd",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PAUSE: {
			Name: "pause",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
		},
		unix.SYS_NANOSLEEP: {
			Name: "nanosleep",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "req",
					Type: argTypeTimespec,
				},
				{
					Name:        "rem",
					Type:        argTypeTimespec,
					Destination: true,
				},
			},
		},
		unix.SYS_GETITIMER: {
			Name: "getitimer",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateWhichTimer,
				},
				{
					Name:        "curr_value",
					Type:        argTypeItimerval,
					Destination: true,
				},
			},
		},
		unix.SYS_ALARM: {
			Name: "alarm",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "seconds",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_SETITIMER: {
			Name: "setitimer",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateWhichTimer,
				},
				{
					Name: "new_value",
					Type: argTypeItimerval,
				},
				{
					Name:      "old_value",
					Type:      argTypeItimerval,
					Annotator: annotation.AnnotateNull,
				},
			},
		},
		unix.SYS_GETPID: {
			Name: "getpid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SENDFILE: {
			Name: "sendfile",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "out_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "in_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_SOCKET: {
			Name: "socket",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "domain",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketDomain,
				},
				{
					Name:      "type",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketType,
				},
				{
					Name:      "protocol",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketProtocol,
				},
			},
		},
		unix.SYS_CONNECT: {
			Name: "connect",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "addr",
					Type:      argTypeSockaddr,
					LenSource: LenSourceNext,
				},
				{
					Name: "addrlen",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_ACCEPT: {
			Name: "accept",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "addr",
					Type:      argTypeSockaddr,
					LenSource: LenSourceNextPointer,
				},
				{
					Name: "addrlen",
					Type: argTypeUnsignedIntPtr,
				},
			},
		},
		unix.SYS_SENDTO: {
			Name: "sendto",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_RECVFROM: {
			Name: "recvfrom",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceNext,
					Destination: true,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
				{
					Name:        "addr",
					Type:        argTypeSockaddr,
					Destination: true,
					LenSource:   LenSourceNextPointer,
				},
				{
					Name:        "addrlen",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
			},
		},
		unix.SYS_SENDMSG: {
			Name: "sendmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "msg",
					Type: argTypeMsghdr,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_RECVMSG: {
			Name: "recvmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "msg",
					Type:        argTypeMsghdr,
					Destination: true,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_SHUTDOWN: {
			Name: "shutdown",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "how",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateShutdownHow,
				},
			},
		},
		unix.SYS_BIND: {
			Name: "bind",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "addr",
					Type: argTypeSockaddr,
				},
				{
					Name: "addrlen",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_LISTEN: {
			Name: "listen",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "backlog",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_GETSOCKNAME: {
			Name: "getsockname",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "addr",
					Type:        argTypeSockaddr,
					Destination: true,
					LenSource:   LenSourceNextPointer,
				},
				{
					Name:        "addrlen",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
			},
		},
		unix.SYS_GETPEERNAME: {
			Name: "getpeername",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "addr",
					Type:        argTypeSockaddr,
					Destination: true,
					LenSource:   LenSourceNextPointer,
				},
				{
					Name:        "addrlen",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
			},
		},
		unix.SYS_SOCKETPAIR: {
			Name: "socketpair",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "domain",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketDomain,
				},
				{
					Name:      "type",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketType,
				},
				{
					Name:      "protocol",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketProtocol,
				},
				{
					Name:       "fds",
					Type:       argTypeIntArray,
					LenSource:  LenSourceFixed,
					FixedCount: 2,
					Annotator: func(arg annotation.Arg, pid int) {
						if underlying, ok := arg.(*Arg); ok {
							for i := 0; i < 2; i++ {
								annotation.AnnotateFd(&underlying.array[i], pid)
							}
						}
					},
				},
			},
		},
		unix.SYS_SETSOCKOPT: {
			Name: "setsockopt",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "level",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketLevel,
				},
				{
					Name:      "optname",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketOption,
				},
				{
					Name:      "optval",
					Type:      argTypeSockoptval,
					LenSource: LenSourceNext,
				},
				{
					Name: "optlen",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_GETSOCKOPT: {
			Name: "getsockopt",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "level",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketLevel,
				},
				{
					Name:      "optname",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSocketOption,
				},
				{
					Name:        "optval",
					Type:        argTypeSockoptval,
					LenSource:   LenSourceNextPointer,
					Destination: true,
				},
				{
					Name:        "optlen",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
			},
		},
		unix.SYS_CLONE: {
			Name: "clone",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fn",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name:      "child_stack",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateCloneFlags,
				},
				{
					Name:      "arg",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
			},
		},
		unix.SYS_FORK: {
			Name: "fork",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_VFORK: {
			Name: "vfork",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_EXECVE: {
			Name: "execve",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name: "argv",
					Type: argTypeStringArray,
				},
				{
					Name: "envp",
					Type: argTypeStringArray,
				},
			},
		},
		unix.SYS_EXIT: {
			Name: "exit",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "status",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_WAIT4: {
			Name: "wait4",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
					/*
						< -1

						wait for any child process whose process group ID is equal to the absolute value of pid.
						-1

						wait for any child process; this is equivalent to calling wait3().
						0

						wait for any child process whose process group ID is equal to that of the calling process.
						> 0

						wait for the child whose process ID is equal to the value of pid.
					*/
				},
				{
					Name:      "status",
					Type:      argTypeWaitStatus,
					Annotator: annotation.AnnotateWaitStatus,
				},
				{
					Name:      "options",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateWaitOptions,
				},
				{
					Name:        "rusage",
					Type:        argTypeRUsage,
					Destination: true,
				},
			},
		},
		unix.SYS_KILL: {
			Name: "kill",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
			},
		},
		unix.SYS_UNAME: {
			Name: "uname",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "buf",
					Type:        argTypeUname,
					Destination: true,
				},
			},
		},
		unix.SYS_SEMGET: {
			Name: "semget",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "key",
					Type: ArgTypeInt,
				},
				{
					Name: "nsems",
					Type: ArgTypeInt,
				},
				{
					Name:      "semflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSemFlags,
				},
			},
		},
		unix.SYS_SEMOP: {
			Name: "semop",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "semid",
					Type: ArgTypeInt,
				},
				{
					Name: "sops",
					Type: argTypeSembuf,
				},
				{
					Name: "nsops",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SEMCTL: {
			Name: "semctl",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "semid",
					Type: ArgTypeInt,
				},
				{
					Name: "semnum",
					Type: ArgTypeInt,
				},
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSemCmd,
				},
				{
					Name:        "arg",
					Type:        ArgTypeAddress,
					Annotator:   annotation.AnnotateNull, // TODO: annotate this based on cmd?
					Destination: true,                    // sometimes true, so always wait
				},
			},
		},
		unix.SYS_SHMDT: {
			Name: "shmdt",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "shmaddr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_MSGGET: {
			Name: "msgget",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "key",
					Type: ArgTypeInt,
				},
				{
					Name:      "msgflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_MSGSND: {
			Name: "msgsnd",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "msqid",
					Type: ArgTypeInt,
				},
				{
					Name:      "msgp",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "msgsz",
					Type: ArgTypeInt,
				},
				{
					Name:      "msgflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_MSGRCV: {
			Name: "msgrcv",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "msqid",
					Type: ArgTypeInt,
				},
				{
					Name:        "msgp",
					Type:        ArgTypeData,
					LenSource:   LenSourceNext,
					Destination: true,
				},
				{
					Name: "msgsz",
					Type: ArgTypeInt,
				},
				{
					Name:      "msgtyp",
					Type:      ArgTypeLong,
					Annotator: annotation.AnnotateMsgType,
				},
				{
					Name:      "msgflg",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_MSGCTL: {
			Name: "msgctl",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "msqid",
					Type: ArgTypeInt,
				},
				{
					Name: "cmd",
					Type: ArgTypeInt,
				},
				{
					Name:        "buf",
					Type:        ArgTypeAddress, // TODO: this needs it's own type so we can print more info
					Destination: true,
				},
			},
		},
		unix.SYS_FCNTL: {
			Name: "fcntl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFcntlCmd,
				},
				{
					Name: "arg",
					Type: ArgTypeAddress,
				},
			},
			Modifier: func(call *Syscall) {

				args := call.args[2:]

				switch call.args[1].raw {
				case unix.F_DUPFD, unix.F_DUPFD_CLOEXEC, unix.F_SETFD, unix.F_SETFL, unix.F_SETOWN, unix.F_SETSIG, unix.F_SETLEASE, unix.F_NOTIFY, unix.F_SETPIPE_SZ:
					args = args[:1]
					args[0] = Arg{
						raw:  args[0].raw,
						t:    ArgTypeInt,
						name: "arg",
					}
				case unix.F_GETLK, unix.F_SETLK, unix.F_SETLKW:
					args = args[:1]
					args[0] = Arg{
						raw:  args[0].raw,
						t:    ArgTypeAddress, // TODO create flock type
						name: "arg",
					}
				default:
					args = nil
				}

				call.args = append(call.args[:2], args...)
			},
		},
		unix.SYS_FLOCK: {
			Name: "flock",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "operation",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFlockOperation,
				},
			},
		},
		unix.SYS_FSYNC: {
			Name: "fsync",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_FDATASYNC: {
			Name: "fdatasync",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_TRUNCATE: {
			Name: "truncate",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "length",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_FTRUNCATE: {
			Name: "ftruncate",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "length",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_GETDENTS: {
			Name: "getdents",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "dirp",
					Type: ArgTypeAddress,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_GETCWD: {
			Name: "getcwd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
					Annotator:   trimToNull,
				},
				{
					Name: "size",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_CHDIR: {
			Name: "chdir",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_FCHDIR: {
			Name: "fchdir",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_RENAME: {
			Name: "rename",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldpath",
					Type: argTypeString,
				},
				{
					Name: "newpath",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_MKDIR: {
			Name: "mkdir",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_RMDIR: {
			Name: "rmdir",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_CREAT: {
			Name: "creat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_LINK: {
			Name: "link",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldpath",
					Type: argTypeString,
				},
				{
					Name: "newpath",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_UNLINK: {
			Name: "unlink",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_SYMLINK: {
			Name: "symlink",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldpath",
					Type: argTypeString,
				},
				{
					Name: "newpath",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_READLINK: {
			Name: "readlink",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceNext,
					Destination: true,
					Annotator:   trimToNull,
				},
				{
					Name: "bufsiz",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CHMOD: {
			Name: "chmod",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_FCHMOD: {
			Name: "fchmod",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_CHOWN: {
			Name: "chown",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "owner",
					Type: ArgTypeInt,
				},
				{
					Name: "group",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FCHOWN: {
			Name: "fchown",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "owner",
					Type: ArgTypeInt,
				},
				{
					Name: "group",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LCHOWN: {
			Name: "lchown",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "owner",
					Type: ArgTypeInt,
				},
				{
					Name: "group",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_UMASK: {
			Name: "umask",
			ReturnValue: ReturnMetadata{
				Type:      ArgTypeInt,
				Annotator: annotation.AnnotateAccMode,
			},
			Args: []ArgMetadata{
				{
					Name:      "mask",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_GETTIMEOFDAY: {
			Name: "gettimeofday",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "tv",
					Type:        argTypeTimeval,
					Destination: true,
				},
				{
					Name:        "tz",
					Type:        argTypeTimezone,
					Destination: true,
				},
			},
		},
		unix.SYS_GETRLIMIT: {
			Name: "getrlimit",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "resource",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRLimitResourceFlags,
				},
				{
					Name:        "rlim",
					Type:        argTypeRLimit,
					Destination: true,
				},
			},
		},
		unix.SYS_GETRUSAGE: {
			Name: "getrusage",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "who",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRUsageWho,
				},
				{
					Name:        "usage",
					Type:        argTypeRUsage,
					Destination: true,
				},
			},
		},
		unix.SYS_SYSINFO: {
			Name: "sysinfo",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "info",
					Type:        argTypeSysinfo,
					Destination: true,
				},
			},
		},
		unix.SYS_TIMES: {
			Name: "times",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "tbuf",
					Type:        argTypeTms,
					Destination: true,
				},
			},
		},
		unix.SYS_PTRACE: {
			Name: "ptrace",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "request",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePtraceRequest,
				},
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "data",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_GETUID: {
			Name: "getuid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SYSLOG: {
			Name: "syslog",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "type",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSyslogType,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					Destination: true,
					LenSource:   LenSourceReturnValue,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETGID: {
			Name: "getgid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SETUID: {
			Name: "setuid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "uid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETGID: {
			Name: "setgid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "gid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETEUID: {
			Name: "geteuid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_GETEGID: {
			Name: "getegid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SETPGID: {
			Name: "setpgid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "pgid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETPPID: {
			Name: "getppid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_GETPGRP: {
			Name: "getpgrp",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SETSID: {
			Name: "setsid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
		},
		unix.SYS_SETREUID: {
			Name: "setreuid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ruid",
					Type: ArgTypeInt,
				},
				{
					Name: "euid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETREGID: {
			Name: "setregid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "rgid",
					Type: ArgTypeInt,
				},
				{
					Name: "egid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETGROUPS: {
			Name: "getgroups",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "size",
					Type:        ArgTypeInt,
					Destination: true,
				},
				{
					Name:        "list",
					Type:        argTypeIntArray,
					Destination: true,
				},
			},
		},
		unix.SYS_SETGROUPS: {
			Name: "setgroups",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name: "list",
					Type: argTypeIntArray,
				},
			},
		},
		unix.SYS_SETRESUID: {
			Name: "setresuid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ruid",
					Type: ArgTypeInt,
				},
				{
					Name: "euid",
					Type: ArgTypeInt,
				},
				{
					Name: "suid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETRESUID: {
			Name: "getresuid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "ruid",
					Type:        ArgTypeInt,
					Destination: true,
				},
				{
					Name:        "euid",
					Type:        ArgTypeInt,
					Destination: true,
				},
				{
					Name:        "suid",
					Type:        ArgTypeInt,
					Destination: true,
				},
			},
		},
		unix.SYS_SETRESGID: {
			Name: "setresgid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "rgid",
					Type: ArgTypeInt,
				},
				{
					Name: "egid",
					Type: ArgTypeInt,
				},
				{
					Name: "sgid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETRESGID: {
			Name: "getresgid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "rgid",
					Type:        ArgTypeInt,
					Destination: true,
				},
				{
					Name:        "egid",
					Type:        ArgTypeInt,
					Destination: true,
				},
				{
					Name:        "sgid",
					Type:        ArgTypeInt,
					Destination: true,
				},
			},
		},
		unix.SYS_GETPGID: {
			Name: "getpgid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETFSUID: {
			Name: "setfsuid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name: "fsuid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETFSGID: {
			Name: "setfsgid",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name: "fsgid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETSID: {
			Name: "getsid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CAPGET: {
			Name: "capget",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "header",
					Type:        argTypeCapUserHeader,
					Destination: true,
				},
				{
					Name:        "data",
					Type:        argTypeCapUserData,
					Destination: true,
				},
			},
		},
		unix.SYS_CAPSET: {
			Name: "capset",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "header",
					Type: argTypeCapUserHeader,
				},
				{
					Name: "data",
					Type: argTypeCapUserData,
				},
			},
		},
		unix.SYS_RT_SIGPENDING: {
			Name: "rt_sigpending",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "set",
					Type: ArgTypeAddress,
				},
				{
					Name: "sigsetsize",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_RT_SIGTIMEDWAIT: {
			Name: "rt_sigtimedwait",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "set",
					Type: ArgTypeAddress,
				},
				{
					Name: "info",
					Type: argTypeSigInfo,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
			},
		},
		unix.SYS_RT_SIGQUEUEINFO: {
			Name: "rt_sigqueueinfo",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
				{
					Name: "info",
					Type: argTypeSigInfo,
				},
			},
		},
		unix.SYS_RT_SIGSUSPEND: {
			Name: "rt_sigsuspend",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name: "set",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_SIGALTSTACK: {
			Name: "sigaltstack",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ss",
					Type: argTypeStack,
				},
				{
					Name:        "old_ss",
					Type:        argTypeStack,
					Destination: true,
				},
			},
		},
		unix.SYS_UTIME: {
			Name: "utime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:        "times",
					Type:        argTypeUtimbuf,
					Destination: true,
				},
			},
		},
		unix.SYS_MKNOD: {
			Name: "mknod",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
				{
					Name:      "dev",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateDevice,
				},
			},
		},
		unix.SYS_USELIB: {
			Name: "uselib",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "library",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_PERSONALITY: {
			Name: "personality",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "persona",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_USTAT: {
			Name: "ustat",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dev",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateDevice,
				},
				{
					Name:        "ubuf",
					Type:        argTypeUstat,
					Destination: true,
				},
			},
		},
		unix.SYS_STATFS: {
			Name: "statfs",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name:        "buf",
					Type:        argTypeStatfs,
					Destination: true,
				},
			},
		},
		unix.SYS_FSTATFS: {
			Name: "fstatfs",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "buf",
					Type:        argTypeStatfs,
					Destination: true,
				},
			},
		},
		unix.SYS_SYSFS: {
			Name: "sysfs",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "option",
					Type: ArgTypeInt,
				},
			},
			Modifier: func(call *Syscall) {
				switch call.args[0].raw {
				case 1: // int sysfs(int option, const char *fsname)
					str, _ := readString(call.pid, call.rawArgs[1])
					call.args = append(call.args, Arg{
						name: "fsname",
						t:    ArgTypeData,
						data: []byte(str),
					})
				case 2: // int sysfs(int option, int fs_index, const char *buf)
					call.args = append(call.args, Arg{
						name: "fs_index",
						t:    ArgTypeInt,
						raw:  call.rawArgs[1],
					})
					str, _ := readString(call.pid, call.rawArgs[2])
					call.args = append(call.args, Arg{
						name: "buf",
						t:    ArgTypeData,
						data: []byte(str),
					})
				}
			},
		},
		unix.SYS_GETPRIORITY: {
			Name: "getpriority",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePriorityWhich,
				},
				{
					Name: "who",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETPRIORITY: {
			Name: "setpriority",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePriorityWhich,
				},
				{
					Name: "who",
					Type: ArgTypeInt,
				},
				{
					Name: "prio",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SCHED_SETPARAM: {
			Name: "sched_setparam",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "param",
					Type: argTypeSchedParam,
				},
			},
		},
		unix.SYS_SCHED_GETPARAM: {
			Name: "sched_getparam",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:        "param",
					Type:        argTypeSchedParam,
					Destination: true,
				},
			},
		},
		unix.SYS_SCHED_SETSCHEDULER: {
			Name: "sched_setscheduler",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "policy",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSchedPolicy,
				},
				{
					Name: "param",
					Type: argTypeSchedParam,
				},
			},
		},
		unix.SYS_SCHED_GETSCHEDULER: {
			Name: "sched_getscheduler",
			ReturnValue: ReturnMetadata{
				Type:      argTypeIntOrErrorCode,
				Annotator: annotation.AnnotateSchedPolicy,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SCHED_GET_PRIORITY_MAX: {
			Name: "sched_get_priority_max",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "policy",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSchedPolicy,
				},
			},
		},
		unix.SYS_SCHED_GET_PRIORITY_MIN: {
			Name: "sched_get_priority_min",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "policy",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSchedPolicy,
				},
			},
		},
		unix.SYS_SCHED_RR_GET_INTERVAL: {
			Name: "sched_rr_get_interval",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:        "interval",
					Type:        argTypeTimespec,
					Destination: true,
				},
			},
		},
		unix.SYS_MLOCK: {
			Name: "mlock",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "addr",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MUNLOCK: {
			Name: "munlock",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "addr",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MLOCKALL: {
			Name: "mlockall",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMlockFlags,
				},
			},
		},
		unix.SYS_MUNLOCKALL: {
			Name: "munlockall",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_VHANGUP: {
			Name: "vhangup",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_MODIFY_LDT: {
			Name: "modify_ldt",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "func",
					Type: ArgTypeInt,
					Annotator: func(arg annotation.Arg, pid int) {
						switch arg.Raw() {
						case 0:
							arg.SetAnnotation("read_ldt", true)
						case 1:
							arg.SetAnnotation("write_ldt", true)
						}
					},
				},
				{
					Name:        "ptr",
					Type:        argTypeUserDesc,
					Destination: true,
				},
				{
					Name: "bytecount",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PIVOT_ROOT: {
			Name: "pivot_root",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "new_root",
					Type: argTypeString,
				},
				{
					Name: "put_old",
					Type: argTypeString,
				},
			},
		},
		unix.SYS__SYSCTL: {
			Name: "sysctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "args",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_PRCTL: {
			Name: "prctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "option",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePrctlOption,
				},
				// TODO: make these options dynamic based on the option via Modifier
				{
					Name: "arg2",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "arg3",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "arg4",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "arg5",
					Type: ArgTypeUnsignedLong,
				},
			},
		},
		unix.SYS_ARCH_PRCTL: {
			Name: "arch_prctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "code",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateArchPrctrlCode,
				},
				{
					Name:        "addr",
					Type:        ArgTypeAddress,
					Destination: true,
				},
			},
		},
		unix.SYS_ADJTIMEX: {
			Name: "adjtimex",
			ReturnValue: ReturnMetadata{
				Type:      argTypeIntOrErrorCode,
				Annotator: annotation.AnnotateClockState,
			},
			Args: []ArgMetadata{
				{
					Name: "buf",
					Type: argTypeTimex,
				},
			},
		},
		unix.SYS_SETRLIMIT: {
			Name: "setrlimit",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "resource",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRLimitResourceFlags,
				},
				{
					Name: "rlim",
					Type: argTypeRLimit,
				},
			},
		},
		unix.SYS_CHROOT: {
			Name: "chroot",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_SYNC: {
			Name: "sync",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_ACCT: {
			Name: "acct",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_SETTIMEOFDAY: {
			Name: "settimeofday",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "tv",
					Type: argTypeTimeval,
				},
				{
					Name: "tz",
					Type: argTypeTimezone,
				},
			},
		},
		unix.SYS_MOUNT: {
			Name: "mount",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "source",
					Type: argTypeString,
				},
				{
					Name: "target",
					Type: argTypeString,
				},
				{
					Name: "filesystemtype",
					Type: argTypeString,
				},
				{
					Name:      "mountflags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateMountFlags,
				},
				{
					Name: "data",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_UMOUNT2: {
			Name: "umount2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "target",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateUmountFlags,
				},
			},
		},
		unix.SYS_SWAPON: {
			Name: "swapon",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_SWAPOFF: {
			Name: "swapoff",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_REBOOT: {
			Name: "reboot",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "magic",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRebootMagic,
				},
				{
					Name:      "magic2",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRebootMagic,
				},
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRebootCmd,
				},
				{
					Name: "arg",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_SETHOSTNAME: {
			Name: "sethostname",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "name",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_SETDOMAINNAME: {
			Name: "setdomainname",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "name",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_IOPL: {
			Name: "iopl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "level",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_IOPERM: {
			Name: "ioperm",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "from",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "num",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "turn_on",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CREATE_MODULE: {
			Name: "create_module",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_INIT_MODULE: {
			Name: "init_module",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "module_image",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "param_values",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_DELETE_MODULE: {
			Name: "delete_module",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateDeleteModuleFlags,
				},
			},
		},
		unix.SYS_GET_KERNEL_SYMS: {
			Name: "get_kernel_syms",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "table",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
			},
		},
		unix.SYS_QUERY_MODULE: {
			Name: "query_module",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateQueryModuleWhich,
				},
				{
					Name: "buf",
					Type: ArgTypeAddress,
				},
				{
					Name: "bufsize",
					Type: ArgTypeInt,
				},
				{
					Name: "ret",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_QUOTACTL: {
			Name: "quotactl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateQuotactlCmd,
				},
				{
					Name: "special",
					Type: argTypeString,
				},
				{
					Name: "id",
					Type: ArgTypeInt,
				},
				{
					Name: "addr",
					Type: ArgTypeAddress, // TODO: annotate structs for each possible cmd
				},
			},
		},
		unix.SYS_NFSSERVCTL: { // removed in 3.1
			Name: "nfsservctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "cmd",
					Type: ArgTypeInt,
				},
				{
					Name: "argp",
					Type: ArgTypeAddress,
				},
				{
					Name: "resp",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_GETPMSG: {
			Name: "getpmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "ctlptr",
					Type: ArgTypeAddress,
				},
				{
					Name: "dataptr",
					Type: ArgTypeAddress,
				},
				{
					Name: "bandp",
					Type: ArgTypeAddress,
				},
				{
					Name: "flagsp",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_PUTPMSG: {
			Name: "putpmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "ctlptr",
					Type: ArgTypeAddress,
				},
				{
					Name: "dataptr",
					Type: ArgTypeAddress,
				},
				{
					Name: "band",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_AFS_SYSCALL: {
			Name: "afs_syscall",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_TUXCALL: {
			Name: "tuxcall",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_SECURITY: {
			Name: "security",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_GETTID: {
			Name: "gettid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_READAHEAD: {
			Name: "readahead",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeLong,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_SETXATTR: {
			Name: "setxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "value",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateXFlags,
				},
			},
		},
		unix.SYS_LSETXATTR: {
			Name: "lsetxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "value",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateXFlags,
				},
			},
		},
		unix.SYS_FSETXATTR: {
			Name: "fsetxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "value",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateXFlags,
				},
			},
		},
		unix.SYS_GETXATTR: {
			Name: "getxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:        "value",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LGETXATTR: {
			Name: "lgetxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:        "value",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FGETXATTR: {
			Name: "fgetxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:        "value",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LISTXATTR: {
			Name: "listxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name:        "list",
					Type:        argTypeStringArray,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LLISTXATTR: {
			Name: "llistxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name:        "list",
					Type:        argTypeStringArray,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FLISTXATTR: {
			Name: "flistxattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "list",
					Type:        argTypeStringArray,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_REMOVEXATTR: {
			Name: "removexattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_LREMOVEXATTR: {
			Name: "lremovexattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_FREMOVEXATTR: {
			Name: "fremovexattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "name",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_TKILL: {
			Name: "tkill",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "tid",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
			},
		},
		unix.SYS_TIME: {
			Name: "time",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "t",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_FUTEX: {
			Name: "futex",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "uaddr",
					Type: ArgTypeAddress,
				},
				{
					Name:      "op",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFutexOp,
				},
				{
					Name: "val",
					Type: ArgTypeInt,
				},
				{
					Name:     "timeout",
					Type:     argTypeTimespec,
					Optional: true,
				},
				{
					Name:     "uaddr2",
					Type:     ArgTypeAddress,
					Optional: true,
				},
				{
					Name:     "val3",
					Type:     ArgTypeInt,
					Optional: true,
				},
			},
		},
		unix.SYS_SCHED_SETAFFINITY: {
			Name: "sched_setaffinity",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name: "user_mask_ptr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_SCHED_GETAFFINITY: {
			Name: "sched_getaffinity",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name: "user_mask_ptr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_SET_THREAD_AREA: {
			Name: "set_thread_area",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "u_info",
					Type: argTypeUserDesc,
				},
			},
		},
		unix.SYS_IO_SETUP: {
			Name: "io_setup",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "nr_events",
					Type: ArgTypeInt,
				},
				{
					Name:        "ctx_idp",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
			},
		},
		unix.SYS_IO_DESTROY: {
			Name: "io_destroy",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ctx_id",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_IO_GETEVENTS: {
			Name: "io_getevents",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ctx_id",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "min_nr",
					Type: ArgTypeLong,
				},
				{
					Name: "nr",
					Type: ArgTypeLong,
				},
				{
					Name: "events",
					Type: argTypeIoEvents,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
			},
		},
		unix.SYS_IO_SUBMIT: {
			Name: "io_submit",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ctx_id",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "nr",
					Type: ArgTypeLong,
				},
				{
					Name: "iocbpp",
					Type: argTypeIoCB,
				},
			},
		},
		unix.SYS_IO_CANCEL: {
			Name: "io_cancel",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ctx_id",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "iocb",
					Type: argTypeIoCB,
				},
				{
					Name: "result",
					Type: argTypeIoEvent,
				},
			},
		},
		unix.SYS_GET_THREAD_AREA: {
			Name: "get_thread_area",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "u_info",
					Type:        argTypeUserDesc,
					Destination: true,
				},
			},
		},
		unix.SYS_LOOKUP_DCOOKIE: {
			Name: "lookup_dcookie",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "cookie64",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "buf",
					Type:      ArgTypeAddress,
					LenSource: LenSourceReturnValue,
					Annotator: trimToNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_EPOLL_CREATE: {
			Name: "epoll_create",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EPOLL_CTL_OLD: {
			Name: "epoll_ctl_old",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name: "op",
					Type: ArgTypeInt,
				},
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "event",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_EPOLL_WAIT_OLD: {
			Name: "epoll_wait_old",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name: "events",
					Type: ArgTypeAddress,
				},
				{
					Name: "maxevents",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_REMAP_FILE_PAGES: {
			Name: "remap_file_pages",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "start",
					Type: ArgTypeAddress,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "prot",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateProt,
				},
				{
					Name: "pgoff",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GETDENTS64: {
			Name: "getdents64",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "dirp",
					Type: ArgTypeAddress,
				},
				{
					Name: "count",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_SET_TID_ADDRESS: {
			Name: "set_tid_address",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "tidptr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_RESTART_SYSCALL: {
			Name: "restart_syscall",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_SEMTIMEDOP: {
			Name: "semtimedop",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "semid",
					Type: ArgTypeInt,
				},
				{
					Name: "sops",
					Type: argTypeSembuf,
				},
				{
					Name: "nsops",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
			},
		},
		unix.SYS_FADVISE64: {
			Name: "fadvise64",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeLong,
				},
				{
					Name: "len",
					Type: ArgTypeLong,
				},
				{
					Name:      "advice",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFAdvice,
				},
			},
		},
		unix.SYS_TIMER_CREATE: {
			Name: "timer_create",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name: "sevp",
					Type: ArgTypeAddress, // TODO: create sigevent type?
				},
				{
					Name: "timerid",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_TIMER_SETTIME: {
			Name: "timer_settime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "timerid",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateTimerFlags,
				},
				{
					Name: "new_value",
					Type: argTypeItimerspec,
				},
				{
					Name:        "old_value",
					Type:        argTypeItimerspec,
					Destination: true,
				},
			},
		},
		unix.SYS_TIMER_GETTIME: {
			Name: "timer_gettime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "timerid",
					Type: ArgTypeInt,
				},
				{
					Name:        "curr_value",
					Type:        argTypeItimerspec,
					Destination: true,
				},
			},
		},
		unix.SYS_TIMER_GETOVERRUN: {
			Name: "timer_getoverrun",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "timerid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_TIMER_DELETE: {
			Name: "timer_delete",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "timerid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CLOCK_SETTIME: {
			Name: "clock_settime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name: "tp",
					Type: argTypeTimespec,
				},
			},
		},
		unix.SYS_CLOCK_GETTIME: {
			Name: "clock_gettime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name:        "tp",
					Type:        argTypeTimespec,
					Destination: true,
				},
			},
		},
		unix.SYS_CLOCK_GETRES: {
			Name: "clock_getres",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "clockid",
					Type: ArgTypeInt,
				},
				{
					Name:        "res",
					Type:        argTypeTimespec,
					Destination: true,
				},
			},
		},
		unix.SYS_CLOCK_NANOSLEEP: {
			Name: "clock_nanosleep",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateTimerFlags,
				},
				{
					Name: "rqtp",
					Type: argTypeTimespec,
				},
				{
					Name:        "rmtp",
					Type:        argTypeTimespec,
					Destination: true,
				},
			},
		},
		unix.SYS_EXIT_GROUP: {
			Name: "exit_group",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "status",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EPOLL_WAIT: {
			Name: "epoll_wait",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name:        "events",
					Type:        argTypeEpollEvent,
					Destination: true,
				},
				{
					Name: "maxevents",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EPOLL_CTL: {
			Name: "epoll_ctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name:      "op",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateEpollCtlOp,
				},
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "event",
					Type: argTypeEpollEvent,
				},
			},
		},
		unix.SYS_TGKILL: {
			Name: "tgkill",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "tgid",
					Type: ArgTypeInt,
				},
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
			},
		},
		unix.SYS_UTIMES: {
			Name: "utimes",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:       "times",
					Type:       argTypeTimevalArray,
					LenSource:  LenSourceFixed,
					FixedCount: 2,
				},
			},
		},
		unix.SYS_VSERVER: {
			Name: "vserver",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_MBIND: {
			Name: "mbind",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "mode",
					Type: ArgTypeInt,
				},
				{
					Name: "nodemask",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "maxnode",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateNumaModeFlag,
				},
			},
		},
		unix.SYS_SET_MEMPOLICY: {
			Name: "set_mempolicy",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateNumaModeFlag,
				},
				{
					Name: "nodemask",
					Type: ArgTypeAddress,
				},
				{
					Name: "maxnode",
					Type: ArgTypeUnsignedLong,
				},
			},
		},
		unix.SYS_GET_MEMPOLICY: {
			Name: "get_mempolicy",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateNumaModeFlag,
				},
				{
					Name: "nodemask",
					Type: ArgTypeAddress,
				},
				{
					Name: "maxnode",
					Type: ArgTypeAddress,
				},
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "flags",
					Type: ArgTypeUnsignedLong, // TODO: annotate flags
				},
			},
		},
		unix.SYS_MQ_OPEN: {
			Name: "mq_open",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "oflag",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
					Optional:  true,
				},
				{
					Name:     "attr",
					Type:     argTypeMqAttr,
					Optional: true,
				},
			},
		},
		unix.SYS_MQ_UNLINK: {
			Name: "mq_unlink",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_MQ_TIMEDSEND: {
			Name: "mq_timedsend",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "mqdes",
					Type: ArgTypeInt,
				},
				{
					Name:      "msg_ptr",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "msg_len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "msg_prio",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "abs_timeout",
					Type: argTypeTimespec,
				},
			},
		},
		unix.SYS_MQ_TIMEDRECEIVE: {
			Name: "mq_timedreceive",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "mqdes",
					Type: ArgTypeInt,
				},
				{
					Name:        "msg_ptr",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "msg_len",
					Type: ArgTypeUnsignedLong,
				},
			},
		},
		unix.SYS_MQ_NOTIFY: {
			Name: "mq_notify",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "mqdes",
					Type: ArgTypeInt,
				},
				{
					Name: "notification",
					Type: ArgTypeAddress, // TODO: sigevent type
				},
			},
		},
		unix.SYS_MQ_GETSETATTR: {
			Name: "mq_getsetattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "mqdes",
					Type: ArgTypeInt,
				},
				{
					Name: "mqstat",
					Type: argTypeMqAttr,
				},
				{
					Name:        "omqstat",
					Type:        argTypeMqAttr,
					Destination: true,
				},
			},
		},
		unix.SYS_KEXEC_LOAD: {
			Name: "kexec_load",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "entry",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "nr_segments",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "segments",
					Type: ArgTypeAddress,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateKexecFlags,
				},
			},
		},
		unix.SYS_WAITID: {
			Name: "waitid",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "idtype",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateIDType,
				},
				{
					Name: "id",
					Type: ArgTypeInt,
				},
				{
					Name:        "infop",
					Type:        argTypeSigInfo,
					Destination: true,
				},
				{
					Name:      "options",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateWaitOptions,
				},
			},
		},
		unix.SYS_ADD_KEY: {
			Name: "add_key",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "type",
					Type: argTypeString,
				},
				{
					Name: "description",
					Type: argTypeString,
				},
				{
					Name:      "payload",
					Type:      ArgTypeData,
					LenSource: LenSourceNext,
				},
				{
					Name: "plen",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "keyring",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateKeyringID,
				},
			},
		},
		unix.SYS_REQUEST_KEY: {
			Name: "request_key",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "type",
					Type: argTypeString,
				},
				{
					Name: "description",
					Type: argTypeString,
				},
				{
					Name: "callout_info",
					Type: argTypeString,
				},
				{
					Name:      "keyring",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateKeyringID,
				},
			},
		},
		unix.SYS_KEYCTL: {
			Name: "keyctl",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateKeyctlCommand,
				},
			},
			// TODO: add a modifier here to dynamically add arguments for specific commands
		},
		unix.SYS_IOPRIO_SET: {
			Name: "ioprio_set",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateIoPrioWhich,
				},
				{
					Name: "who",
					Type: ArgTypeInt,
				},
				{
					Name: "ioprio",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_IOPRIO_GET: {
			Name: "ioprio_get",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "which",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateIoPrioWhich,
				},
				{
					Name: "who",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_INOTIFY_INIT: {
			Name: "inotify_init",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
		},
		unix.SYS_INOTIFY_ADD_WATCH: {
			Name: "inotify_add_watch",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "mask",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_INOTIFY_RM_WATCH: {
			Name: "inotify_rm_watch",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "wd",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MIGRATE_PAGES: {
			Name: "migrate_pages",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "maxnode",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "old_nodes",
					Type: ArgTypeAddress,
				},
				{
					Name: "new_nodes",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_OPENAT: {
			Name: "openat",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Optional:  true,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_MKDIRAT: {
			Name: "mkdirat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
			},
		},
		unix.SYS_MKNODAT: {
			Name: "mknodat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
				{
					Name:      "dev",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateDevice,
				},
			},
		},
		unix.SYS_FCHOWNAT: {
			Name: "fchownat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name: "owner",
					Type: ArgTypeInt,
				},
				{
					Name: "group",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_FUTIMESAT: {
			Name: "futimesat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:       "times",
					Type:       argTypeTimevalArray,
					LenSource:  LenSourceFixed,
					FixedCount: 2,
				},
			},
		},
		unix.SYS_NEWFSTATAT: {
			Name: "newfstatat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name:        "statbuf",
					Type:        argTypeStat,
					Destination: true,
				},
				{
					Name:      "flag",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_UNLINKAT: {
			Name: "unlinkat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_RENAMEAT: {
			Name: "renameat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "olddfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "oldname",
					Type: argTypeString,
				},
				{
					Name:      "newdfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "newname",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_LINKAT: {
			Name: "linkat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "olddfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "oldname",
					Type: argTypeString,
				},
				{
					Name:      "newdfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "newname",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_SYMLINKAT: {
			Name: "symlinkat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldname",
					Type: argTypeString,
				},
				{
					Name:      "newdfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "newname",
					Type: argTypeString,
				},
			},
		},
		unix.SYS_READLINKAT: {
			Name: "readlinkat",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					LenSource: LenSourceReturnValue,
				},
				{
					Name: "bufsiz",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FCHMODAT: {
			Name: "fchmodat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_FACCESSAT: {
			Name: "faccessat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
					Optional:  true,
				},
			},
		},
		unix.SYS_PSELECT6: {
			Name: "pselect6",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "nfds",
					Type: ArgTypeInt,
				},
				{
					Name: "readfds",
					Type: argTypeFdSet,
				},
				{
					Name: "writefds",
					Type: argTypeFdSet,
				},
				{
					Name: "exceptfds",
					Type: argTypeFdSet,
				},
				{
					Name: "timeout",
					Type: argTypeTimeval,
				},
				{
					Name: "sigmask",
					Type: ArgTypeAddress, // TODO: sigset type
				},
			},
		},
		unix.SYS_PPOLL: {
			Name: "ppoll",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "fds",
					Type: argTypePollFdArray,
				},
				{
					Name: "nfds",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
				{
					Name: "sigmask",
					Type: ArgTypeAddress, // TODO: sigset type (check all sigmask properties)
				},
			},
		},
		unix.SYS_UNSHARE: {
			Name: "unshare",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateCloneFlags,
				},
			},
		},
		unix.SYS_SET_ROBUST_LIST: {
			Name: "set_robust_list",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "head",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_GET_ROBUST_LIST: {
			Name: "get_robust_list",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "head_ptr",
					Type: ArgTypeAddress,
				},
				{
					Name: "len_ptr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_SPLICE: {
			Name: "splice",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd_in",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "off_in",
					Type: ArgTypeAddress,
				},
				{
					Name:      "fd_out",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "off_out",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSpliceFlags,
				},
			},
		},
		unix.SYS_TEE: {
			Name: "tee",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd_in",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "fd_out",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSpliceFlags,
				},
			},
		},
		unix.SYS_SYNC_FILE_RANGE: {
			Name: "sync_file_range",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name: "nbytes",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSyncFileRangeFlags,
				},
			},
		},
		unix.SYS_VMSPLICE: {
			Name: "vmsplice",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "iov",
					Type:        argTypeIovecArray,
					Destination: true,
				},
				{
					Name: "nr_segs",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSpliceFlags,
				},
			},
		},
		unix.SYS_MOVE_PAGES: {
			Name: "move_pages",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "nr_pages",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "pages",
					Type: ArgTypeAddress,
				},
				{
					Name: "nodes",
					Type: argTypeIntArray,
				},
				{
					Name:        "status",
					Type:        argTypeIntArray,
					Destination: true,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateNumaModeFlag,
				},
			},
		},
		unix.SYS_UTIMENSAT: {
			Name: "utimensat",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dirfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:       "times",
					Type:       argTypeTimespecArray,
					LenSource:  LenSourceFixed,
					FixedCount: 2,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_EPOLL_PWAIT: {
			Name: "epoll_pwait",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name:        "events",
					Type:        argTypeEpollEvent,
					Destination: true,
				},
				{
					Name: "maxevents",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: ArgTypeInt,
				},
				{
					Name: "sigmask",
					Type: ArgTypeAddress, // TODO: sigset type
				},
			},
		},
		unix.SYS_SIGNALFD: {
			Name: "signalfd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "mask",
					Type: ArgTypeAddress, // TODO: sigset type
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignalFdFlags,
				},
			},
		},
		unix.SYS_TIMERFD_CREATE: {
			Name: "timerfd_create",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateTimerFdFlags,
				},
			},
		},
		unix.SYS_EVENTFD: {
			Name: "eventfd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "count",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateEventFdFlags,
				},
			},
		},
		unix.SYS_FALLOCATE: {
			Name: "fallocate",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFallocateMode,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_TIMERFD_SETTIME: {
			Name: "timerfd_settime",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateTimerFdFlags,
				},
				{
					Name: "utmr",
					Type: argTypeItimerspec,
				},
				{
					Name:        "otmr",
					Type:        argTypeItimerspec,
					Destination: true,
				},
			},
		},
		unix.SYS_TIMERFD_GETTIME: {
			Name: "timerfd_gettime",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "otmr",
					Type:        argTypeItimerspec,
					Destination: true,
				},
			},
		},
		unix.SYS_ACCEPT4: {
			Name: "accept4",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "addr",
					Type:      argTypeSockaddr,
					LenSource: LenSourceNextPointer,
				},
				{
					Name: "addrlen",
					Type: argTypeUnsignedIntPtr,
				},
			},
		},
		unix.SYS_SIGNALFD4: {
			Name: "signalfd4",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "mask",
					Type: ArgTypeAddress, // TODO: sigset type
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignalFdFlags,
				},
			},
		},
		unix.SYS_EVENTFD2: {
			Name: "eventfd2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "count",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EPOLL_CREATE1: {
			Name: "epoll_create1",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_DUP3: {
			Name: "dup3",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "oldfd",
					Type: ArgTypeInt,
				},
				{
					Name: "newfd",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PIPE2: {
			Name: "pipe2",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "pipefd",
					Type:        argTypeIntArray,
					LenSource:   LenSourceFixed,
					FixedCount:  2,
					Destination: true,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
					// TODO: annotate pipe2 flags
				},
			},
		},
		unix.SYS_INOTIFY_INIT1: {
			Name: "inotify_init1",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PREADV: {
			Name: "preadv",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "iov",
					Type:        argTypeIovecArray,
					Destination: true,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PWRITEV: {
			Name: "pwritev",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_RT_TGSIGQUEUEINFO: {
			Name: "rt_tgsigqueueinfo",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "tgid",
					Type: ArgTypeInt,
				},
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
				{
					Name: "uinfo",
					Type: argTypeSigInfo,
				},
			},
		},
		unix.SYS_PERF_EVENT_OPEN: {
			Name: "perf_event_open",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "attr",
					Type: ArgTypeAddress,
				},
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "cpu",
					Type: ArgTypeInt,
				},
				{
					Name: "group_fd",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotatePerfFlags,
				},
			},
		},
		unix.SYS_RECVMMSG: {
			Name: "recvmmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "msgvec",
					Type:        argTypeMMsgHdrArray,
					Destination: true,
				},
				{
					Name: "vlen",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
				{
					Name: "timeout",
					Type: argTypeTimeval,
				},
			},
		},
		unix.SYS_FANOTIFY_INIT: {
			Name: "fanotify_init",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateFANotifyFlags,
				},
				{
					Name:      "event_f_flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateFANotifyEventFlags,
				},
			},
		},
		unix.SYS_FANOTIFY_MARK: {
			Name: "fanotify_mark",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fanotify_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotation.AnnotateFANotifyMarkFlags,
				},
				{
					Name: "mask",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_PRLIMIT64: {
			Name: "prlimit64",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:      "resource",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRLimitResourceFlags,
				},
				{
					Name: "new_rlim",
					Type: argTypeRLimit,
				},
				{
					Name:        "old_rlim",
					Type:        argTypeRLimit,
					Destination: true,
				},
			},
		},
		unix.SYS_NAME_TO_HANDLE_AT: {
			Name: "name_to_handle_at",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "handle",
					Type: ArgTypeAddress, // TODO: create type for this
				},
				{
					Name:        "mount_id",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
			},
		},
		unix.SYS_OPEN_BY_HANDLE_AT: {
			Name: "open_by_handle_at",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "mount_fd",
					Type: ArgTypeInt,
				},
				{
					Name: "handle",
					Type: ArgTypeAddress, // TODO: create type for this
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
			},
		},
		unix.SYS_CLOCK_ADJTIME: {
			Name: "clock_adjtime",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "clk_id",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
				{
					Name: "utx",
					Type: argTypeTimex,
				},
			},
		},
		unix.SYS_SYNCFS: {
			Name: "syncfs",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
			},
		},
		unix.SYS_SENDMMSG: {
			Name: "sendmmsg",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "msgvec",
					Type: argTypeMMsgHdrArray,
				},
				{
					Name: "vlen",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMsgFlags,
				},
			},
		},
		unix.SYS_SETNS: {
			Name: "setns",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "nstype",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateCloneFlags,
				},
			},
		},
		unix.SYS_GETCPU: {
			Name: "getcpu",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "cpu",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
				{
					Name:        "node",
					Type:        argTypeUnsignedIntPtr,
					Destination: true,
				},
				{
					Name:      "tcache",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
			},
		},
		unix.SYS_PROCESS_VM_READV: {
			Name: "process_vm_readv",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "local_iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "liovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "remote_iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "riovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PROCESS_VM_WRITEV: {
			Name: "process_vm_writev",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "local_iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "liovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "remote_iov",
					Type: argTypeIovecArray,
				},
				{
					Name: "riovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_KCMP: {
			Name: "kcmp",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid1",
					Type: ArgTypeInt,
				},
				{
					Name: "pid2",
					Type: ArgTypeInt,
				},
				{
					Name:      "type",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateKcmpType,
				},
				{
					Name: "idx1",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "idx2",
					Type: ArgTypeUnsignedLong,
				},
			},
		},
		unix.SYS_FINIT_MODULE: {
			Name: "finit_module",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "uargs",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateModuleInitFlags,
				},
			},
		},
		unix.SYS_SCHED_SETATTR: {
			Name: "sched_setattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "attr",
					Type: argTypeSchedAttr,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSchedPolicy,
				},
			},
		},
		unix.SYS_SCHED_GETATTR: {
			Name: "sched_getattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name:        "attr",
					Type:        argTypeSchedAttr,
					Destination: true,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSchedPolicy,
				},
			},
		},
		unix.SYS_RENAMEAT2: {
			Name: "renameat2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "oldfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "oldname",
					Type: argTypeString,
				},
				{
					Name:      "newfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "newname",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_SECCOMP: {
			Name: "seccomp",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "op",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSeccompOp,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
				{
					Name: "uargs",
					Type: ArgTypeAddress,
				},
				// TODO: add a modifier to annotate flags and args depending on op
			},
		},
		unix.SYS_GETRANDOM: {
			Name: "getrandom",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "buf",
					Type:        ArgTypeData,
					LenSource:   LenSourceReturnValue,
					Destination: true,
				},
				{
					Name: "count",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateRandomFlags,
				},
			},
		},
		unix.SYS_MEMFD_CREATE: {
			Name: "memfd_create",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "name",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMemfdFlags,
				},
			},
		},
		unix.SYS_KEXEC_FILE_LOAD: {
			Name: "kexec_file_load",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "kernel_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "initrd_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "cmdline_len",
					Type: ArgTypeInt,
				},
				{
					Name: "cmdline_ptr",
					Type: ArgTypeAddress,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateKexecFlags,
				},
			},
		},
		unix.SYS_BPF: {
			Name: "bpf",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateBPFCmd,
				},
				{
					Name: "attr",
					Type: ArgTypeAddress,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EXECVEAT: {
			Name: "execveat",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dirfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "argv",
					Type: argTypeStringArray,
				},
				{
					Name: "envp",
					Type: argTypeStringArray,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
			},
		},
		unix.SYS_USERFAULTFD: {
			Name: "userfaultfd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateOpenFlags,
				},
			},
		},
		unix.SYS_MEMBARRIER: {
			Name: "membarrier",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "cmd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMembarrierCmd,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMembarrierCmdFlags,
				},
				{
					Name: "cpuid",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MLOCK2: {
			Name: "mlock2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "addr",
					Type:      ArgTypeAddress,
					Annotator: annotation.AnnotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMlockFlags,
				},
			},
		},
		unix.SYS_COPY_FILE_RANGE: {
			Name: "copy_file_range",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd_in",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "off_in",
					Type: argTypeUnsignedInt64Ptr,
				},
				{
					Name:      "fd_out",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "off_out",
					Type: argTypeUnsignedInt64Ptr,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PREADV2: {
			Name: "preadv2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "iov",
					Type:        argTypeIovecArray,
					Destination: true,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePReadWrite2Flags,
				},
			},
		},
		unix.SYS_PWRITEV2: {
			Name: "pwritev2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:        "iov",
					Type:        argTypeIovecArray,
					Destination: true,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePReadWrite2Flags,
				},
			},
		},
		unix.SYS_PKEY_MPROTECT: {
			Name: "pkey_mprotect",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeInt,
				},
				{
					Name:      "prot",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateProt,
				},
				{
					Name: "pkey",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PKEY_ALLOC: {
			Name: "pkey_alloc",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
				{
					Name:      "access_rights",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotatePkeyAccessRights,
				},
			},
		},
		unix.SYS_PKEY_FREE: {
			Name: "pkey_free",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pkey",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_STATX: {
			Name: "statx",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dirfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
				{
					Name:      "mask",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateStatxMask,
				},
				{
					Name:        "statxbuf",
					Type:        argTypeStatX,
					Destination: true,
				},
			},
		},
		unix.SYS_IO_PGETEVENTS: {
			Name: "io_pgetevents",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "ctx_id",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "nr",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:        "events",
					Type:        argTypeIoEvents,
					Destination: true,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
				{
					Name: "usig",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_RSEQ: {
			Name: "rseq",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "rseq",
					Type: ArgTypeAddress,
				},
				{
					Name: "rseq_len",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
			},
		},
		unix.SYS_PIDFD_SEND_SIGNAL: {
			Name: "pidfd_send_signal",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pidfd",
					Type: ArgTypeInt,
				},
				{
					Name:      "sig",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateSignal,
				},
				{
					Name: "info",
					Type: argTypeSigInfo,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_IO_URING_SETUP: {
			Name: "io_uring_setup",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "entries",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name:        "params",
					Type:        argTypeIoUringParams,
					Destination: true,
				},
			},
		},
		unix.SYS_IO_URING_ENTER: {
			Name: "io_uring_enter",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "to_submit",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name: "min_complete",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotation.AnnotateIoUringEnterFlags,
				},
				{
					Name: "sig",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_IO_URING_REGISTER: {
			Name: "io_uring_register",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "opcode",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotation.AnnotateIORingOpCode,
				},
				{
					Name: "arg",
					Type: ArgTypeAddress,
				},
				{
					Name: "nr_args",
					Type: ArgTypeUnsignedInt,
				},
			},
		},
		unix.SYS_OPEN_TREE: {
			Name: "open_tree",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "filename",
					Type: argTypeString,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MOVE_MOUNT: {
			Name: "move_mount",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "from_dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "from_pathname",
					Type: argTypeString,
				},
				{
					Name:      "to_dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "to_pathname",
					Type: argTypeString,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FSOPEN: {
			Name: "fsopen",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "fs_name",
					Type: argTypeString,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FSCONFIG: {
			Name: "fsconfig",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fs_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "cmd",
					Type: ArgTypeInt,
				},
				{
					Name: "key",
					Type: argTypeString,
				},
				{
					Name: "value",
					Type: argTypeString,
				},
				{
					Name: "aux",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FSMOUNT: {
			Name: "fsmount",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fs_fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
				{
					Name: "ms_flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FSPICK: {
			Name: "fspick",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PIDFD_OPEN: {
			Name: "pidfd_open",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pid",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CLONE3: {
			Name: "clone3",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "clone_args",
					Type: argTypeCloneArgs,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_CLOSE_RANGE: {
			Name: "close_range",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "first",
					Type: ArgTypeInt,
				},
				{
					Name: "last",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateCloseRangeFlags,
				},
			},
		},
		unix.SYS_OPENAT2: {
			Name: "openat2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dirfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name: "how",
					Type: argTypeOpenHow,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_PIDFD_GETFD: {
			Name: "pidfd_getfd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "pidfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name:      "targetfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FACCESSAT2: {
			Name: "faccessat2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "pathname",
					Type: argTypeString,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAccMode,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
					Optional:  true,
				},
			},
		},
		unix.SYS_PROCESS_MADVISE: {
			Name: "process_madvise",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "pidfd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "iovec",
					Type: argTypeIovecArray,
				},
				{
					Name: "iovcnt",
					Type: ArgTypeInt,
				},
				{
					Name:      "advice",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMAdviseAdvice,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EPOLL_PWAIT2: {
			Name: "epoll_pwait2",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "epfd",
					Type: ArgTypeInt,
				},
				{
					Name:        "events",
					Type:        argTypeEpollEvent,
					Destination: true,
				},
				{
					Name: "maxevents",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
				{
					Name: "sigmask",
					Type: ArgTypeAddress, // TODO: sigset type
				},
			},
		},
		unix.SYS_MOUNT_SETATTR: {
			Name: "mount_setattr",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "path",
					Type: argTypeString,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateAtFlags,
				},
				{
					Name: "uattr",
					Type: argTypeMountAttr,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_QUOTACTL_FD: {
			Name: "quotactl_fd",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateFd,
				},
				{
					Name: "cmd",
					Type: ArgTypeInt,
				},
				{
					Name: "id",
					Type: ArgTypeInt,
				},
				{
					Name: "addr",
					Type: ArgTypeAddress,
				},
			},
		},
		unix.SYS_LANDLOCK_CREATE_RULESET: {
			Name: "landlock_create_ruleset",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ruleset_attr",
					Type: argTypeLandlockRulesetAttr,
				},
				{
					Name: "size",
					Type: ArgTypeInt,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateLandlockFlags,
				},
			},
		},
		unix.SYS_LANDLOCK_ADD_RULE: {
			Name: "landlock_add_rule",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ruleset_fd",
					Type: ArgTypeInt,
				},
				{
					Name:      "rule_type",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateLandlockRuleType,
				},
				{
					Name: "rule_attr",
					Type: ArgTypeAddress,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_LANDLOCK_RESTRICT_SELF: {
			Name: "landlock_restrict_self",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "ruleset_fd",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_MEMFD_SECRET: {
			Name: "memfd_secret",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateMemfdFlags,
				},
			},
		},
		unix.SYS_PROCESS_MRELEASE: {
			Name: "process_mrelease",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "pidfd",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_FUTEX_WAITV: {
			Name: "futex_waitv",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "waiters",
					Type: ArgTypeAddress,
				},
				{
					Name: "nr_waiters",
					Type: ArgTypeInt,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
				{
					Name: "timeout",
					Type: argTypeTimespec,
				},
				{
					Name:      "clockid",
					Type:      ArgTypeInt,
					Annotator: annotation.AnnotateClockID,
				},
			},
		},
		unix.SYS_SET_MEMPOLICY_HOME_NODE: {
			Name: "set_mempolicy_home_node",
			ReturnValue: ReturnMetadata{
				Type: argTypeIntOrErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name: "start",
					Type: ArgTypeAddress,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "home_node",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name: "flags",
					Type: ArgTypeInt,
				},
			},
		},
	}
)

func trimToNull(arg annotation.Arg, _ int) {
	// clean up terminating null byte from output
	underlying := arg.(*Arg)
	if index := bytes.Index(underlying.data, []byte{0}); index >= 0 {
		underlying.data = underlying.data[:index]
	}
}

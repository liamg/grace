//go:build amd64

package tracer

import (
	"strings"
	"syscall"

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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					CountFrom:   CountLocationResult,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					CountFrom: CountLocationNext,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "filename",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotateOpenFlags,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Annotator: annotateAccMode,
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
					Annotator: annotateFd,
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
					Name:      "filename",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:        "stat",
					Type:        ArgTypeStat,
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
					Annotator: annotateFd,
				},
				{
					Name:        "stat",
					Type:        ArgTypeStat,
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
					Name:      "filename",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:        "stat",
					Type:        ArgTypeStat,
					Destination: true,
				},
			},
		},
		unix.SYS_POLL: {
			Name: "poll",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "ufds",
					Type:        ArgTypePollFdArray,
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
					Annotator: annotateFd,
				},
				{
					Name: "offset",
					Type: ArgTypeInt,
				},
				{
					Name:      "whence",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotateWhence,
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
					Annotator: annotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedLong,
				},
				{
					Name:      "prot",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotateProt,
				},
				{
					Name: "flags",
					Type: ArgTypeUnsignedLong,
					Annotator: func(arg *Arg, _ int) {
						var joins []string

						switch arg.Raw() & 0x3 {
						case unix.MAP_SHARED:
							joins = append(joins, "MAP_SHARED")
						case unix.MAP_PRIVATE:
							joins = append(joins, "MAP_PRIVATE")
						case unix.MAP_SHARED_VALIDATE:
							joins = append(joins, "MAP_SHARED_VALIDATE")
						}

						mapConsts := map[int]string{
							unix.MAP_32BIT:            "MAP_32BIT",
							unix.MAP_ANONYMOUS:        "MAP_ANONYMOUS",
							unix.MAP_DENYWRITE:        "MAP_DENYWRITE",
							unix.MAP_EXECUTABLE:       "MAP_EXECUTABLE",
							unix.MAP_FILE:             "MAP_FILE",
							unix.MAP_FIXED:            "MAP_FIXED",
							unix.MAP_FIXED_NOREPLACE:  "MAP_FIXED_NOREPLACE",
							unix.MAP_GROWSDOWN:        "MAP_GROWSDOWN",
							unix.MAP_HUGETLB:          "MAP_HUGETLB",
							21 << unix.MAP_HUGE_SHIFT: "MAP_HUGE_2MB",
							30 << unix.MAP_HUGE_SHIFT: "MAP_HUGE_1GB",
							unix.MAP_LOCKED:           "MAP_LOCKED",
							unix.MAP_NONBLOCK:         "MAP_NONBLOCK",
							unix.MAP_NORESERVE:        "MAP_NORESERVE",
							unix.MAP_POPULATE:         "MAP_POPULATE",
							unix.MAP_STACK:            "MAP_STACK",
							unix.MAP_SYNC:             "MAP_SYNC",
						}

						for flag, str := range mapConsts {
							if arg.Raw()&uintptr(flag) > 0 {
								joins = append(joins, str)
							}
						}
						arg.annotation = strings.Join(joins, "|")
						arg.replace = arg.annotation != ""
					},
				},
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
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
					Annotator: annotateNull,
				},
				{
					Name: "len",
					Type: ArgTypeUnsignedInt,
				},
				{
					Name:      "prot",
					Type:      ArgTypeUnsignedLong,
					Annotator: annotateProt,
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
					Annotator: annotateNull,
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
					Annotator: annotateNull,
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
					Type: ArgTypeSigAction,
				},
				{
					Name:        "oldact",
					Type:        ArgTypeSigAction,
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
					Annotator: annotateSigProcMaskFlags,
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
					Annotator: annotateFd,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{

				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:        "buf",
					Type:        ArgTypeData,
					CountFrom:   CountLocationNext,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{

				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:      "buf",
					Type:      ArgTypeData,
					CountFrom: CountLocationNext,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name: "iov",
					Type: ArgTypeIovecArray,
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
				Type: ArgTypeInt,
			},
			Args: []ArgMetadata{
				{
					Name:      "fd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name: "iov",
					Type: ArgTypeIovecArray,
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
					Name:      "path",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:      "mode",
					Type:      ArgTypeUnsignedInt,
					Annotator: annotateAccMode,
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
					Type:        ArgTypeIntArray,
					CountFrom:   CountLocationFixed,
					FixedCount:  2,
					Destination: true,
				},
			},
		},
		// <-- progress
		// TODO: add ReturnValue to everything below here...
		unix.SYS_PIPE2: {
			Name: "pipe2",
			ReturnValue: ReturnMetadata{
				Type: ArgTypeErrorCode,
			},
			Args: []ArgMetadata{
				{
					Name:        "pipefd",
					Type:        ArgTypeIntArray,
					CountFrom:   CountLocationFixed,
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
		unix.SYS_OPENAT: {
			Name: "openat",
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:      "filename",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:      "flags",
					Type:      ArgTypeInt,
					Annotator: annotateOpenFlags,
				},
				{
					Name:      "mode",
					Type:      ArgTypeInt,
					Optional:  true,
					Annotator: annotateAccMode,
				},
			},
		},
		unix.SYS_NEWFSTATAT: {
			Name: "newfstatat",
			Args: []ArgMetadata{
				{
					Name:      "dfd",
					Type:      ArgTypeInt,
					Annotator: annotateFd,
				},
				{
					Name:      "filename",
					Type:      ArgTypeData,
					CountFrom: CountLocationNullTerminator,
				},
				{
					Name:        "statbuf",
					Type:        ArgTypeStat,
					Destination: true,
				},
				{
					Name: "flag",
					Type: ArgTypeInt, // TODO: annotate flags
				},
			},
		},
		unix.SYS_EXIT: {
			Name: "exit",
			Args: []ArgMetadata{
				{
					Name: "status",
					Type: ArgTypeInt,
				},
			},
		},
		unix.SYS_EXIT_GROUP: {
			Name: "exit_group",
			Args: []ArgMetadata{
				{
					Name: "status",
					Type: ArgTypeInt,
				},
			},
		},
	}
)

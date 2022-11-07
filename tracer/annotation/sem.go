package annotation

/*
#include <linux/ipc.h>
#include <linux/sem.h>
*/
import "C"

import (
	"strings"

	"golang.org/x/sys/unix"
)

var semFlags = map[int]string{
	unix.IPC_CREAT:  "IPC_CREAT",
	unix.IPC_EXCL:   "IPC_EXCL",
	unix.IPC_NOWAIT: "IPC_NOWAIT",
	unix.S_IRUSR:    "S_IRUSR",
	unix.S_IWUSR:    "S_IWUSR",
	unix.S_IXUSR:    "S_IXUSR",
	unix.S_IRGRP:    "S_IRGRP",
	unix.S_IWGRP:    "S_IWGRP",
	unix.S_IXGRP:    "S_IXGRP",
	unix.S_IROTH:    "S_IROTH",
	unix.S_IWOTH:    "S_IWOTH",
	unix.S_IXOTH:    "S_IXOTH",
}

func AnnotateSemFlags(arg Arg, pid int) {
	var joins []string
	for flag, name := range semFlags {
		if arg.Raw()&uintptr(flag) != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var semCmds = map[int]string{
	unix.IPC_STAT:  "IPC_STAT",
	unix.IPC_SET:   "IPC_SET",
	unix.IPC_RMID:  "IPC_RMID",
	C.IPC_INFO:     "IPC_INFO",
	C.SEM_INFO:     "SEM_INFO",
	C.SEM_STAT:     "SEM_STAT",
	C.SEM_STAT_ANY: "SEM_STAT_ANY",
	C.GETALL:       "GETALL",
	C.GETNCNT:      "GETNCNT",
	C.GETPID:       "GETPID",
	C.GETVAL:       "GETVAL",
	C.GETZCNT:      "GETZCNT",
	C.SETALL:       "SETALL",
	C.SETVAL:       "SETVAL",
}

func AnnotateSemCmd(arg Arg, _ int) {
	if name, ok := semCmds[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

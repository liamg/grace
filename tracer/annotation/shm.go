package annotation

/*
#include <linux/shm.h>
*/
import "C"
import (
	"strings"

	"golang.org/x/sys/unix"
)

var shmCommands = map[int]string{
	unix.IPC_STAT: "IPC_STAT",
	unix.IPC_SET:  "IPC_SET",
	unix.IPC_RMID: "IPC_RMID",
}

func AnnotateSHMCTLCommand(arg Arg, _ int) {
	str, ok := shmCommands[int(arg.Raw())]
	arg.SetAnnotation(str, ok)
}

var shmAtFlags = map[int]string{
	unix.SHM_RDONLY: "SHM_RDONLY",
	C.SHM_REMAP:     "SHM_REMAP",
}

func AnnotateSHMAtFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range shmAtFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var shmGetFlags = map[int]string{
	unix.IPC_CREAT:  "IPC_CREAT",
	unix.IPC_EXCL:   "IPC_EXCL",
	C.SHM_HUGETLB:   "SHM_HUGETLB",
	C.SHM_HUGE_1GB:  "SHM_HUGE_1GB",
	C.SHM_HUGE_2MB:  "SHM_HUGE_2MB",
	C.SHM_NORESERVE: "SHM_NO_RESERVE",
}

func AnnotateSHMGetFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range shmGetFlags {
		if arg.Raw()&uintptr(flag) != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var fcntlCmds = map[int]string{
	unix.F_DUPFD:         "F_DUPFD",
	unix.F_DUPFD_CLOEXEC: "F_DUPFD_CLOEXEC",
	unix.F_GETFD:         "F_GETFD",
	unix.F_SETFD:         "F_SETFD",
	unix.F_GETFL:         "F_GETFL",
	unix.F_SETFL:         "F_SETFL",
	unix.F_GETLK:         "F_GETLK",
	unix.F_SETLK:         "F_SETLK",
	unix.F_SETLKW:        "F_SETLKW",
	unix.F_SETOWN:        "F_SETOWN",
	unix.F_GETOWN:        "F_GETOWN",
	unix.F_GETOWN_EX:     "F_GETOWN_EX",
	unix.F_SETOWN_EX:     "F_SETOWN_EX",
	unix.F_GETSIG:        "F_GETSIG",
	unix.F_SETSIG:        "F_SETSIG",
	unix.F_GETLEASE:      "F_GETLEASE",
	unix.F_SETLEASE:      "F_SETLEASE",
	unix.F_NOTIFY:        "F_NOTIFY",
	unix.F_GETPIPE_SZ:    "F_GETPIPE_SZ",
	unix.F_SETPIPE_SZ:    "F_SETPIPE_SZ",
}

func AnnotateFcntlCmd(arg Arg, pid int) {
	if str, ok := fcntlCmds[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var atFlags = map[int]string{
	unix.AT_SYMLINK_NOFOLLOW: "AT_SYMLINK_NOFOLLOW",
	unix.AT_REMOVEDIR:        "AT_REMOVEDIR",
	unix.AT_SYMLINK_FOLLOW:   "AT_SYMLINK_FOLLOW",
	unix.AT_EMPTY_PATH:       "AT_EMPTY_PATH",
	unix.AT_NO_AUTOMOUNT:     "AT_NO_AUTOMOUNT",
}

func AnnotateAtFlags(arg Arg, pid int) {
	var joins []string
	for flag, str := range atFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

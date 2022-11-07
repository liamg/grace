package annotation

import (
	"golang.org/x/sys/unix"
)

var rlimitFlags = map[int]string{
	unix.RLIMIT_AS:         "RLIMIT_AS",
	unix.RLIMIT_CORE:       "RLIMIT_CORE",
	unix.RLIMIT_CPU:        "RLIMIT_CPU",
	unix.RLIMIT_DATA:       "RLIMIT_DATA",
	unix.RLIMIT_FSIZE:      "RLIMIT_FSIZE",
	unix.RLIMIT_LOCKS:      "RLIMIT_LOCKS",
	unix.RLIMIT_MEMLOCK:    "RLIMIT_MEMLOCK",
	unix.RLIMIT_MSGQUEUE:   "RLIMIT_MSGQUEUE",
	unix.RLIMIT_NICE:       "RLIMIT_NICE",
	unix.RLIMIT_NOFILE:     "RLIMIT_NOFILE",
	unix.RLIMIT_NPROC:      "RLIMIT_NPROC",
	unix.RLIMIT_RSS:        "RLIMIT_RSS",
	unix.RLIMIT_RTPRIO:     "RLIMIT_RTPRIO",
	unix.RLIMIT_RTTIME:     "RLIMIT_RTTIME",
	unix.RLIMIT_SIGPENDING: "RLIMIT_SIGPENDING",
	unix.RLIMIT_STACK:      "RLIMIT_STACK",
}

func AnnotateRLimitResourceFlags(arg Arg, pid int) {
	if str, ok := rlimitFlags[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

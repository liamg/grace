package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var perfFlags = map[int]string{
	unix.PERF_FLAG_FD_CLOEXEC:  "PERF_FLAG_FD_CLOEXEC",
	unix.PERF_FLAG_FD_OUTPUT:   "PERF_FLAG_FD_OUTPUT",
	unix.PERF_FLAG_FD_NO_GROUP: "PERF_FLAG_FD_NO_GROUP",
	unix.PERF_FLAG_PID_CGROUP:  "PERF_FLAG_PID_CGROUP",
}

func AnnotatePerfFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range perfFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

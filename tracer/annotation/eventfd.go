package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var eventFdFlags = map[int]string{
	unix.O_CLOEXEC:     "EFD_CLOEXEC",
	unix.O_NONBLOCK:    "EFD_NONBLOCK",
	unix.EFD_SEMAPHORE: "EFD_SEMAPHORE",
}

func AnnotateEventFdFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range eventFdFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

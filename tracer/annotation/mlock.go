package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var mlockFlags = map[int]string{
	unix.MCL_CURRENT: "MCL_CURRENT",
	unix.MCL_FUTURE:  "MCL_FUTURE",
	unix.MCL_ONFAULT: "MCL_ONFAULT",
}

func AnnotateMlockFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range mlockFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), true)
}

package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var p2Flags = map[int]string{
	unix.RWF_APPEND: "RWF_APPEND",
	unix.RWF_DSYNC:  "RWF_DSYNC",
	unix.RWF_HIPRI:  "RWF_HIPRI",
	unix.RWF_NOWAIT: "RWF_NOWAIT",
	unix.RWF_SYNC:   "RWF_SYNC",
}

func AnnotatePReadWrite2Flags(arg Arg, _ int) {
	var joins []string
	for flag, str := range p2Flags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

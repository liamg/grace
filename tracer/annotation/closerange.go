package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var closeRangeFlags = map[int]string{
	unix.CLOSE_RANGE_CLOEXEC: "CLOSE_RANGE_CLOEXEC",
	unix.CLOSE_RANGE_UNSHARE: "CLOSE_RANGE_UNSHARE",
}

func AnnotateCloseRangeFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range closeRangeFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var spliceFlags = map[int]string{
	unix.SPLICE_F_MOVE:     "SPLICE_F_MOVE",
	unix.SPLICE_F_NONBLOCK: "SPLICE_F_NONBLOCK",
	unix.SPLICE_F_MORE:     "SPLICE_F_MORE",
	unix.SPLICE_F_GIFT:     "SPLICE_F_GIFT",
}

func AnnotateSpliceFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range spliceFlags {
		if (int(arg.Raw())&flag != 0) || (int(arg.Raw()) == flag) {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

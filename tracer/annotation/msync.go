package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var mSyncFlags = map[int]string{
	unix.MS_ASYNC:      "MS_ASYNC",
	unix.MS_INVALIDATE: "MS_INVALIDATE",
	unix.MS_SYNC:       "MS_SYNC",
}

func AnnotateMSyncFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range mSyncFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

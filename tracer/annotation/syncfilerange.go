package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var syncFileRangeFlags = map[int]string{
	unix.SYNC_FILE_RANGE_WAIT_BEFORE: "SYNC_FILE_RANGE_WAIT_BEFORE",
	unix.SYNC_FILE_RANGE_WRITE:       "SYNC_FILE_RANGE_WRITE",
	unix.SYNC_FILE_RANGE_WAIT_AFTER:  "SYNC_FILE_RANGE_WAIT_AFTER",
}

func AnnotateSyncFileRangeFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range syncFileRangeFlags {
		if (int(arg.Raw())&flag != 0) || (int(arg.Raw()) == flag) {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

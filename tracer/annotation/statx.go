package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var statxFlags = map[int]string{
	unix.STATX_TYPE:        "STATX_TYPE",
	unix.STATX_MODE:        "STATX_MODE",
	unix.STATX_NLINK:       "STATX_NLINK",
	unix.STATX_UID:         "STATX_UID",
	unix.STATX_GID:         "STATX_GID",
	unix.STATX_ATIME:       "STATX_ATIME",
	unix.STATX_MTIME:       "STATX_MTIME",
	unix.STATX_CTIME:       "STATX_CTIME",
	unix.STATX_INO:         "STATX_INO",
	unix.STATX_SIZE:        "STATX_SIZE",
	unix.STATX_BLOCKS:      "STATX_BLOCKS",
	unix.STATX_BASIC_STATS: "STATX_BASIC_STATS",
	unix.STATX_BTIME:       "STATX_BTIME",
	unix.STATX_MNT_ID:      "STATX_MNT_ID",
	unix.STATX_ALL:         "STATX_ALL",
}

func AnnotateStatxMask(arg Arg, _ int) {
	var joins []string
	for flag, str := range statxFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

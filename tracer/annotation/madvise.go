package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

func AnnotateMAdviseAdvice(arg Arg, _ int) {
	flags := map[int]string{
		unix.MADV_NORMAL:      "MADV_NORMAL",
		unix.MADV_RANDOM:      "MADV_RANDOM",
		unix.MADV_SEQUENTIAL:  "MADV_SEQUENTIAL",
		unix.MADV_WILLNEED:    "MADV_WILLNEED",
		unix.MADV_DONTNEED:    "MADV_DONTNEED",
		unix.MADV_REMOVE:      "MADV_REMOVE",
		unix.MADV_DONTFORK:    "MADV_DONTFORK",
		unix.MADV_DOFORK:      "MADV_DOFORK",
		unix.MADV_MERGEABLE:   "MADV_MERGEABLE",
		unix.MADV_UNMERGEABLE: "MADV_UNMERGEABLE",
		unix.MADV_HUGEPAGE:    "MADV_HUGEPAGE",
		unix.MADV_NOHUGEPAGE:  "MADV_NOHUGEPAGE",
		unix.MADV_DONTDUMP:    "MADV_DONTDUMP",
		unix.MADV_DODUMP:      "MADV_DODUMP",
		unix.MADV_HWPOISON:    "MADV_HWPOISON",
		unix.MADV_COLD:        "MADV_COLD",
		unix.MADV_PAGEOUT:     "MADV_PAGEOUT",
	}
	var joins []string
	for flag, str := range flags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

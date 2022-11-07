package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var memfdFlags = map[int]string{
	unix.MFD_CLOEXEC:       "MFD_CLOEXEC",
	unix.MFD_ALLOW_SEALING: "MFD_ALLOW_SEALING",
	unix.MFD_HUGETLB:       "MFD_HUGETLB",
	unix.MFD_HUGE_16GB:     "MFD_HUGE_16GB",
	unix.MFD_HUGE_16MB:     "MFD_HUGE_16MB",
	unix.MFD_HUGE_1GB:      "MFD_HUGE_1GB",
	unix.MFD_HUGE_1MB:      "MFD_HUGE_1MB",
	unix.MFD_HUGE_256MB:    "MFD_HUGE_256MB",
	unix.MFD_HUGE_2GB:      "MFD_HUGE_2GB",
	unix.MFD_HUGE_2MB:      "MFD_HUGE_2MB",
	unix.MFD_HUGE_32MB:     "MFD_HUGE_32MB",
	unix.MFD_HUGE_512KB:    "MFD_HUGE_512KB",
	unix.MFD_HUGE_512MB:    "MFD_HUGE_512MB",
	unix.MFD_HUGE_64KB:     "MFD_HUGE_64KB",
	unix.MFD_HUGE_8MB:      "MFD_HUGE_8MB",
}

func AnnotateMemfdFlags(arg Arg, pid int) {
	var joins []string
	for flag, name := range memfdFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

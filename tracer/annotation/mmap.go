package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

func AnnotateMMapFlags(arg Arg, _ int) {
	var joins []string

	switch arg.Raw() & 0x3 {
	case unix.MAP_SHARED:
		joins = append(joins, "MAP_SHARED")
	case unix.MAP_PRIVATE:
		joins = append(joins, "MAP_PRIVATE")
	case unix.MAP_SHARED_VALIDATE:
		joins = append(joins, "MAP_SHARED_VALIDATE")
	}

	mapConsts := map[int]string{
		unix.MAP_32BIT:            "MAP_32BIT",
		unix.MAP_ANONYMOUS:        "MAP_ANONYMOUS",
		unix.MAP_DENYWRITE:        "MAP_DENYWRITE",
		unix.MAP_EXECUTABLE:       "MAP_EXECUTABLE",
		unix.MAP_FILE:             "MAP_FILE",
		unix.MAP_FIXED:            "MAP_FIXED",
		unix.MAP_FIXED_NOREPLACE:  "MAP_FIXED_NOREPLACE",
		unix.MAP_GROWSDOWN:        "MAP_GROWSDOWN",
		unix.MAP_HUGETLB:          "MAP_HUGETLB",
		21 << unix.MAP_HUGE_SHIFT: "MAP_HUGE_2MB",
		30 << unix.MAP_HUGE_SHIFT: "MAP_HUGE_1GB",
		unix.MAP_LOCKED:           "MAP_LOCKED",
		unix.MAP_NONBLOCK:         "MAP_NONBLOCK",
		unix.MAP_NORESERVE:        "MAP_NORESERVE",
		unix.MAP_POPULATE:         "MAP_POPULATE",
		unix.MAP_STACK:            "MAP_STACK",
		unix.MAP_SYNC:             "MAP_SYNC",
	}

	for flag, str := range mapConsts {
		if arg.Raw()&uintptr(flag) > 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

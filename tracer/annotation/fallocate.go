package annotation

import "golang.org/x/sys/unix"

var fallocateModes = map[int]string{
	unix.FALLOC_FL_KEEP_SIZE:      "FALLOC_FL_KEEP_SIZE",
	unix.FALLOC_FL_PUNCH_HOLE:     "FALLOC_FL_PUNCH_HOLE",
	unix.FALLOC_FL_COLLAPSE_RANGE: "FALLOC_FL_COLLAPSE_RANGE",
	unix.FALLOC_FL_ZERO_RANGE:     "FALLOC_FL_ZERO_RANGE",
	unix.FALLOC_FL_INSERT_RANGE:   "FALLOC_FL_INSERT_RANGE",
	unix.FALLOC_FL_UNSHARE_RANGE:  "FALLOC_FL_UNSHARE_RANGE",
}

func AnnotateFallocateMode(arg Arg, _ int) {
	if str, ok := fallocateModes[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

package annotation

import "golang.org/x/sys/unix"

var flockOps = map[int]string{
	unix.LOCK_SH: "LOCK_SH",
	unix.LOCK_EX: "LOCK_EX",
	unix.LOCK_NB: "LOCK_NB",
	unix.LOCK_UN: "LOCK_UN",
}

func AnnotateFlockOperation(arg Arg, _ int) {
	if str, ok := flockOps[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

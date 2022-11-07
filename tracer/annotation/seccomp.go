package annotation

/*
#include <linux/seccomp.h>
*/
import "C"

var seccompOps = map[int]string{
	C.SECCOMP_SET_MODE_STRICT:  "SECCOMP_SET_MODE_STRICT",
	C.SECCOMP_SET_MODE_FILTER:  "SECCOMP_SET_MODE_FILTER",
	C.SECCOMP_GET_ACTION_AVAIL: "SECCOMP_GET_ACTION_AVAIL",
}

func AnnotateSeccompOp(arg Arg, pid int) {
	if s, ok := seccompOps[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	}
}

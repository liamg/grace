package annotation

/*
#include <fcntl.h>
*/
import "C"

var fadviceFlags = map[int]string{
	C.POSIX_FADV_NORMAL:     "POSIX_FADV_NORMAL",
	C.POSIX_FADV_RANDOM:     "POSIX_FADV_RANDOM",
	C.POSIX_FADV_SEQUENTIAL: "POSIX_FADV_SEQUENTIAL",
	C.POSIX_FADV_WILLNEED:   "POSIX_FADV_WILLNEED",
	C.POSIX_FADV_DONTNEED:   "POSIX_FADV_DONTNEED",
	C.POSIX_FADV_NOREUSE:    "POSIX_FADV_NOREUSE",
}

func AnnotateFAdvice(arg Arg, pid int) {
	if str, ok := fadviceFlags[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

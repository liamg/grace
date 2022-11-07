package annotation

import "golang.org/x/sys/unix"

var shutdownHow = map[int]string{
	unix.SHUT_RD:   "SHUT_RD",
	unix.SHUT_WR:   "SHUT_WR",
	unix.SHUT_RDWR: "SHUT_RDWR",
}

func AnnotateShutdownHow(arg Arg, _ int) {
	str, ok := shutdownHow[int(arg.Raw())]
	arg.SetAnnotation(str, ok)
}

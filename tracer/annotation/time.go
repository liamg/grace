package annotation

import "golang.org/x/sys/unix"

var clockStates = map[int]string{
	unix.TIME_OK:    "TIME_OK",
	unix.TIME_INS:   "TIME_INS",
	unix.TIME_DEL:   "TIME_DEL",
	unix.TIME_OOP:   "TIME_OOP",
	unix.TIME_WAIT:  "TIME_WAIT",
	unix.TIME_ERROR: "TIME_ERROR",
}

func AnnotateClockState(arg Arg, _ int) {
	if s, ok := clockStates[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	}
}

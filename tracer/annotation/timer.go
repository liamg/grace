package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

func AnnotateWhichTimer(arg Arg, _ int) {
	switch arg.Raw() {
	case unix.ITIMER_REAL:
		arg.SetAnnotation("ITIMER_REAL", true)
	case unix.ITIMER_VIRTUAL:
		arg.SetAnnotation("ITIMER_VIRTUAL", true)
	case unix.ITIMER_PROF:
		arg.SetAnnotation("ITIMER_PROF", true)
	}
}

var timerFlags = map[int]string{
	unix.TIMER_ABSTIME: "TIMER_ABSTIME",
}

func AnnotateTimerFlags(arg Arg, _ int) {
	if str, ok := timerFlags[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var timerFdFlags = map[int]string{
	unix.O_CLOEXEC:               "TFD_CLOEXEC",
	unix.O_NONBLOCK:              "TFD_NONBLOCK",
	unix.TFD_TIMER_ABSTIME:       "TFD_TIMER_ABSTIME",
	unix.TFD_TIMER_CANCEL_ON_SET: "TFD_TIMER_CANCEL_ON_SET",
}

func AnnotateTimerFdFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range timerFdFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

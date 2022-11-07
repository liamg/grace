package annotation

import "golang.org/x/sys/unix"

var clockIds = map[int]string{
	unix.CLOCK_REALTIME:           "CLOCK_REALTIME",
	unix.CLOCK_MONOTONIC:          "CLOCK_MONOTONIC",
	unix.CLOCK_PROCESS_CPUTIME_ID: "CLOCK_PROCESS_CPUTIME_ID",
	unix.CLOCK_THREAD_CPUTIME_ID:  "CLOCK_THREAD_CPUTIME_ID",
	unix.CLOCK_MONOTONIC_RAW:      "CLOCK_MONOTONIC_RAW",
	unix.CLOCK_REALTIME_COARSE:    "CLOCK_REALTIME_COARSE",
	unix.CLOCK_MONOTONIC_COARSE:   "CLOCK_MONOTONIC_COARSE",
	unix.CLOCK_BOOTTIME:           "CLOCK_BOOTTIME",
	unix.CLOCK_REALTIME_ALARM:     "CLOCK_REALTIME_ALARM",
	unix.CLOCK_BOOTTIME_ALARM:     "CLOCK_BOOTTIME_ALARM",
}

func AnnotateClockID(arg Arg, pid int) {
	if str, ok := clockIds[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

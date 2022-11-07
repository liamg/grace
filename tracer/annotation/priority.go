package annotation

import "golang.org/x/sys/unix"

var priorityWhiches = map[int]string{
	unix.PRIO_PROCESS: "PRIO_PROCESS",
	unix.PRIO_PGRP:    "PRIO_PGRP",
	unix.PRIO_USER:    "PRIO_USER",
}

func AnnotatePriorityWhich(arg Arg, _ int) {
	if s, ok := priorityWhiches[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	}
}

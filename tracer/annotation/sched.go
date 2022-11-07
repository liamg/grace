package annotation

/*
#include <linux/sched.h>
*/
import "C"

var schedPolicies = map[int]string{
	C.SCHED_NORMAL: "SCHED_NORMAL",
	C.SCHED_FIFO:   "SCHED_FIFO",
	C.SCHED_RR:     "SCHED_RR",
	C.SCHED_BATCH:  "SCHED_BATCH",
	C.SCHED_IDLE:   "SCHED_IDLE",
}

func AnnotateSchedPolicy(arg Arg, pid int) {
	if s, ok := schedPolicies[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	}
}

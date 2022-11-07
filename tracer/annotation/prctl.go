package annotation

import "golang.org/x/sys/unix"

var prctlOptions = map[int]string{
	unix.PR_CAP_AMBIENT:              "PR_CAP_AMBIENT",
	unix.PR_CAP_AMBIENT_RAISE:        "PR_CAP_AMBIENT_RAISE",
	unix.PR_CAP_AMBIENT_LOWER:        "PR_CAP_AMBIENT_LOWER",
	unix.PR_CAP_AMBIENT_IS_SET:       "PR_CAP_AMBIENT_IS_SET",
	unix.PR_CAP_AMBIENT_CLEAR_ALL:    "PR_CAP_AMBIENT_CLEAR_ALL",
	unix.PR_CAPBSET_READ:             "PR_CAPBSET_READ",
	unix.PR_CAPBSET_DROP:             "PR_CAPBSET_DROP",
	unix.PR_SET_CHILD_SUBREAPER:      "PR_SET_CHILD_SUBREAPER",
	unix.PR_GET_CHILD_SUBREAPER:      "PR_GET_CHILD_SUBREAPER",
	unix.PR_SET_ENDIAN:               "PR_SET_ENDIAN",
	unix.PR_GET_ENDIAN:               "PR_GET_ENDIAN",
	unix.PR_SET_FPEMU:                "PR_SET_FPEMU",
	unix.PR_GET_FPEMU:                "PR_GET_FPEMU",
	unix.PR_SET_FPEXC:                "PR_SET_FPEXC",
	unix.PR_GET_FPEXC:                "PR_GET_FPEXC",
	unix.PR_SET_FP_MODE:              "PR_SET_FP_MODE",
	unix.PR_GET_FP_MODE:              "PR_GET_FP_MODE",
	unix.PR_MPX_ENABLE_MANAGEMENT:    "PR_MPX_ENABLE_MANAGEMENT",
	unix.PR_MPX_DISABLE_MANAGEMENT:   "PR_MPX_DISABLE_MANAGEMENT",
	unix.PR_SET_KEEPCAPS:             "PR_SET_KEEPCAPS",
	unix.PR_GET_KEEPCAPS:             "PR_GET_KEEPCAPS",
	unix.PR_MCE_KILL:                 "PR_MCE_KILL",
	unix.PR_MCE_KILL_GET:             "PR_MCE_KILL_GET",
	unix.PR_SET_MM:                   "PR_SET_MM",
	unix.PR_SET_MM_START_STACK:       "PR_SET_MM_START_STACK",
	unix.PR_SET_MM_START_BRK:         "PR_SET_MM_START_BRK",
	unix.PR_SET_MM_EXE_FILE:          "PR_SET_MM_EXE_FILE",
	unix.PR_SET_MM_MAP:               "PR_SET_MM_MAP",
	unix.PR_SET_NAME:                 "PR_SET_NAME",
	unix.PR_GET_NAME:                 "PR_GET_NAME",
	unix.PR_SET_NO_NEW_PRIVS:         "PR_SET_NO_NEW_PRIVS",
	unix.PR_GET_NO_NEW_PRIVS:         "PR_GET_NO_NEW_PRIVS",
	unix.PR_SET_PTRACER:              "PR_SET_PTRACER",
	unix.PR_SET_SECCOMP:              "PR_SET_SECCOMP",
	unix.PR_GET_SECCOMP:              "PR_GET_SECCOMP",
	unix.PR_SET_SECUREBITS:           "PR_SET_SECUREBITS",
	unix.PR_GET_SECUREBITS:           "PR_GET_SECUREBITS",
	unix.PR_SET_THP_DISABLE:          "PR_SET_THP_DISABLE",
	unix.PR_TASK_PERF_EVENTS_DISABLE: "PR_TASK_PERF_EVENTS_DISABLE",
	unix.PR_TASK_PERF_EVENTS_ENABLE:  "PR_TASK_PERF_EVENTS_ENABLE",
	unix.PR_GET_THP_DISABLE:          "PR_GET_THP_DISABLE",
	unix.PR_GET_TID_ADDRESS:          "PR_GET_TID_ADDRESS",
	unix.PR_SET_TIMERSLACK:           "PR_SET_TIMERSLACK",
	unix.PR_GET_TIMERSLACK:           "PR_GET_TIMERSLACK",
	unix.PR_SET_TSC:                  "PR_SET_TSC",
	unix.PR_GET_TSC:                  "PR_GET_TSC",
	unix.PR_MCE_KILL_CLEAR:           "PR_MCE_KILL_CLEAR",
	unix.PR_SET_VMA:                  "PR_SET_VMA",
}

func AnnotatePrctlOption(arg Arg, _ int) {
	if s, ok := prctlOptions[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	}
}

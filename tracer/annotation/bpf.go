package annotation

import "golang.org/x/sys/unix"

var bpfCmds = map[int]string{
	unix.BPF_MAP_CREATE:                  "BPF_MAP_CREATE",
	unix.BPF_MAP_LOOKUP_ELEM:             "BPF_MAP_LOOKUP_ELEM",
	unix.BPF_MAP_UPDATE_ELEM:             "BPF_MAP_UPDATE_ELEM",
	unix.BPF_MAP_DELETE_ELEM:             "BPF_MAP_DELETE_ELEM",
	unix.BPF_MAP_GET_NEXT_KEY:            "BPF_MAP_GET_NEXT_KEY",
	unix.BPF_PROG_LOAD:                   "BPF_PROG_LOAD",
	unix.BPF_OBJ_PIN:                     "BPF_OBJ_PIN",
	unix.BPF_OBJ_GET:                     "BPF_OBJ_GET",
	unix.BPF_PROG_ATTACH:                 "BPF_PROG_ATTACH",
	unix.BPF_PROG_DETACH:                 "BPF_PROG_DETACH",
	unix.BPF_PROG_TEST_RUN:               "BPF_PROG_TEST_RUN",
	unix.BPF_PROG_GET_NEXT_ID:            "BPF_PROG_GET_NEXT_ID",
	unix.BPF_MAP_GET_NEXT_ID:             "BPF_MAP_GET_NEXT_ID",
	unix.BPF_PROG_GET_FD_BY_ID:           "BPF_PROG_GET_FD_BY_ID",
	unix.BPF_MAP_GET_FD_BY_ID:            "BPF_MAP_GET_FD_BY_ID",
	unix.BPF_OBJ_GET_INFO_BY_FD:          "BPF_OBJ_GET_INFO_BY_FD",
	unix.BPF_PROG_QUERY:                  "BPF_PROG_QUERY",
	unix.BPF_RAW_TRACEPOINT_OPEN:         "BPF_RAW_TRACEPOINT_OPEN",
	unix.BPF_BTF_LOAD:                    "BPF_BTF_LOAD",
	unix.BPF_BTF_GET_FD_BY_ID:            "BPF_BTF_GET_FD_BY_ID",
	unix.BPF_TASK_FD_QUERY:               "BPF_TASK_FD_QUERY",
	unix.BPF_MAP_LOOKUP_AND_DELETE_ELEM:  "BPF_MAP_LOOKUP_AND_DELETE_ELEM",
	unix.BPF_MAP_FREEZE:                  "BPF_MAP_FREEZE",
	unix.BPF_BTF_GET_NEXT_ID:             "BPF_BTF_GET_NEXT_ID",
	unix.BPF_MAP_LOOKUP_BATCH:            "BPF_MAP_LOOKUP_BATCH",
	unix.BPF_MAP_LOOKUP_AND_DELETE_BATCH: "BPF_MAP_LOOKUP_AND_DELETE_BATCH",
	unix.BPF_MAP_UPDATE_BATCH:            "BPF_MAP_UPDATE_BATCH",
	unix.BPF_MAP_DELETE_BATCH:            "BPF_MAP_DELETE_BATCH",
	unix.BPF_LINK_CREATE:                 "BPF_LINK_CREATE",
	unix.BPF_LINK_UPDATE:                 "BPF_LINK_UPDATE",
	unix.BPF_LINK_GET_FD_BY_ID:           "BPF_LINK_GET_FD_BY_ID",
	unix.BPF_LINK_GET_NEXT_ID:            "BPF_LINK_GET_NEXT_ID",
	unix.BPF_ENABLE_STATS:                "BPF_ENABLE_STATS",
	unix.BPF_ITER_CREATE:                 "BPF_ITER_CREATE",
	unix.BPF_LINK_DETACH:                 "BPF_LINK_DETACH",
}

func AnnotateBPFCmd(arg Arg, _ int) {
	if v, ok := bpfCmds[int(arg.Raw())]; ok {
		arg.SetAnnotation(v, true)
	}
}

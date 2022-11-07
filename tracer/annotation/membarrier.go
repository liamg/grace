package annotation

import "strings"

/*
#include <linux/membarrier.h>
*/
import "C"

var membarrierCmds = map[int]string{
	C.MEMBARRIER_CMD_QUERY:                                "MEMBARRIER_CMD_QUERY",
	C.MEMBARRIER_CMD_GLOBAL:                               "MEMBARRIER_CMD_GLOBAL",
	C.MEMBARRIER_CMD_REGISTER_GLOBAL_EXPEDITED:            "MEMBARRIER_CMD_REGISTER_GLOBAL_EXPEDITED",
	C.MEMBARRIER_CMD_PRIVATE_EXPEDITED:                    "MEMBARRIER_CMD_PRIVATE_EXPEDITED",
	C.MEMBARRIER_CMD_REGISTER_PRIVATE_EXPEDITED:           "MEMBARRIER_CMD_REGISTER_PRIVATE_EXPEDITED",
	C.MEMBARRIER_CMD_PRIVATE_EXPEDITED_SYNC_CORE:          "MEMBARRIER_CMD_PRIVATE_EXPEDITED_SYNC_CORE",
	C.MEMBARRIER_CMD_REGISTER_PRIVATE_EXPEDITED_SYNC_CORE: "MEMBARRIER_CMD_REGISTER_PRIVATE_EXPEDITED_SYNC_CORE",
}

var membarrierCmdFlags = map[int]string{
	1: "MEMBARRIER_CMD_FLAG_CPU",
}

func AnnotateMembarrierCmd(arg Arg, _ int) {
	if name, ok := membarrierCmds[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

func AnnotateMembarrierCmdFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range membarrierCmdFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

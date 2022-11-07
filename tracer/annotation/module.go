package annotation

import "C"
import (
	"strings"

	"golang.org/x/sys/unix"
)

var deleteModuleFlags = map[int]string{
	unix.O_NONBLOCK: "O_NONBLOCK",
	unix.O_TRUNC:    "O_TRUNC",
}

func AnnotateDeleteModuleFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range deleteModuleFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var queryModuleWhiches = map[int]string{
	1: "QM_MODULES",
	2: "QM_DEPS",
	3: "QM_REFS",
	4: "QM_SYMBOLS",
	5: "QM_INFO",
	6: "QM_EXPORTS",
}

func AnnotateQueryModuleWhich(arg Arg, _ int) {
	if name, ok := queryModuleWhiches[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

var moduleInitFlags = map[int]string{
	unix.MODULE_INIT_IGNORE_MODVERSIONS: "MODULE_INIT_IGNORE_MODVERSIONS",
	unix.MODULE_INIT_IGNORE_VERMAGIC:    "MODULE_INIT_IGNORE_VERMAGIC",
}

func AnnotateModuleInitFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range moduleInitFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var cloneFlags = map[int]string{
	unix.CLONE_CHILD_CLEARTID: "CLONE_CHILD_CLEARTID",
	unix.CLONE_CHILD_SETTID:   "CLONE_CHILD_SETTID",
	unix.CLONE_FILES:          "CLONE_FILES",
	unix.CLONE_FS:             "CLONE_FS",
	unix.CLONE_IO:             "CLONE_IO",
	unix.CLONE_NEWCGROUP:      "CLONE_NEWCGROUP",
	unix.CLONE_NEWIPC:         "CLONE_NEWIPC",
	unix.CLONE_NEWNET:         "CLONE_NEWNET",
	unix.CLONE_NEWNS:          "CLONE_NEWNS",
	unix.CLONE_NEWPID:         "CLONE_NEWPID",
	unix.CLONE_NEWUSER:        "CLONE_NEWUSER",
	unix.CLONE_NEWUTS:         "CLONE_NEWUTS",
	unix.CLONE_PARENT:         "CLONE_PARENT",
	unix.CLONE_PARENT_SETTID:  "CLONE_PARENT_SETTID",
	unix.CLONE_PTRACE:         "CLONE_PTRACE",
	unix.CLONE_SETTLS:         "CLONE_SETTLS",
	unix.CLONE_SIGHAND:        "CLONE_SIGHAND",
	unix.CLONE_SYSVSEM:        "CLONE_SYSVSEM",
	unix.CLONE_THREAD:         "CLONE_THREAD",
	unix.CLONE_UNTRACED:       "CLONE_UNTRACED",
	unix.CLONE_VFORK:          "CLONE_VFORK",
	unix.CLONE_VM:             "CLONE_VM",
}

func AnnotateCloneFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range cloneFlags {
		if arg.Raw()&uintptr(flag) != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

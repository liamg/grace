package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var kexecLoadFlags = map[int]string{
	unix.KEXEC_FILE_UNLOAD:       "KEXEC_FILE_UNLOAD",
	unix.KEXEC_FILE_ON_CRASH:     "KEXEC_FILE_ON_CRASH",
	unix.KEXEC_FILE_NO_INITRAMFS: "KEXEC_FILE_NO_INITRAMFS",
}

func AnnotateKexecFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range kexecLoadFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

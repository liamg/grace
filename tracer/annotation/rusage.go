package annotation

import "golang.org/x/sys/unix"

var rusageWho = map[int]string{
	unix.RUSAGE_CHILDREN: "RUSAGE_CHILDREN",
	unix.RUSAGE_SELF:     "RUSAGE_SELF",
	unix.RUSAGE_THREAD:   "RUSAGE_THREAD",
}

func AnnotateRUsageWho(arg Arg, _ int) {
	if name, ok := rusageWho[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

package annotation

var ioPrioWhiches = map[int]string{
	1: "IOPRIO_WHO_PROCESS",
	2: "IOPRIO_WHO_PGRP",
	3: "IOPRIO_WHO_USER",
}

func AnnotateIoPrioWhich(arg Arg, _ int) {
	if which, ok := ioPrioWhiches[int(arg.Raw())]; ok {
		arg.SetAnnotation(which, true)
	}
}

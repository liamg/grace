package annotation

var kcmpTypes = map[int]string{
	0: "KCMP_FILE",
	1: "KCMP_VM",
	2: "KCMP_FILES",
	3: "KCMP_FS",
	4: "KCMP_SIGHAND",
	5: "KCMP_IO",
	6: "KCMP_SYSVSEM",
	7: "KCMP_TYPES",
	8: "KCMP_EPOLL_TFD",
}

func AnnotateKcmpType(arg Arg, pid int) {
	if str, ok := kcmpTypes[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

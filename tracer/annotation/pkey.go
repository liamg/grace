package annotation

var pkeyAccessRights = map[int]string{
	1: "PKEY_DISABLE_ACCESS",
	2: "PKEY_DISABLE_WRITE",
}

func AnnotatePkeyAccessRights(arg Arg, _ int) {
	if name, ok := pkeyAccessRights[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

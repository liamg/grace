package annotation

func AnnotateNull(arg Arg, _ int) {
	if arg.Raw() == 0 {
		arg.SetAnnotation("NULL", true)
	}
}

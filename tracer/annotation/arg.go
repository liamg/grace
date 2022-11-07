package annotation

type Arg interface {
	Raw() uintptr
	SetAnnotation(annotation string, replace bool)
}

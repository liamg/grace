package annotation

import "golang.org/x/sys/unix"

func AnnotateWhence(arg Arg, _ int) {
	switch int(arg.Raw()) {
	case unix.SEEK_SET:
		arg.SetAnnotation("SEEK_SET", true)
	case unix.SEEK_CUR:
		arg.SetAnnotation("SEEK_CUR", true)
	case unix.SEEK_END:
		arg.SetAnnotation("SEEK_END", true)
	case unix.SEEK_DATA:
		arg.SetAnnotation("SEEK_DATA", true)
	case unix.SEEK_HOLE:
		arg.SetAnnotation("SEEK_HOLE", true)
	}
}

package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var randomFlags = map[int]string{
	unix.GRND_RANDOM:   "GRND_RANDOM",
	unix.GRND_NONBLOCK: "GRND_NONBLOCK",
}

func AnnotateRandomFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range randomFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

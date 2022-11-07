package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var xFlags = map[int]string{
	unix.XATTR_CREATE:  "XATTR_CREATE",
	unix.XATTR_REPLACE: "XATTR_REPLACE",
}

func AnnotateXFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range xFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

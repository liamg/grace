package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

func AnnotateAccMode(arg Arg, _ int) {
	var joins []string
	if arg.Raw() == unix.F_OK {
		joins = append(joins, "F_OK")
	} else {
		if arg.Raw()&unix.R_OK > 0 {
			joins = append(joins, "R_OK")
		}
		if arg.Raw()&unix.W_OK > 0 {
			joins = append(joins, "W_OK")
		}
		if arg.Raw()&unix.X_OK > 0 {
			joins = append(joins, "X_OK")
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

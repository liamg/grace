package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

func AnnotateProt(arg Arg, _ int) {
	if arg.Raw() == unix.PROT_NONE {
		arg.SetAnnotation("PROT_NONE", true)
		return
	}
	var joins []string
	if arg.Raw()&unix.PROT_READ > 0 {
		joins = append(joins, "PROT_READ")
	}
	if arg.Raw()&unix.PROT_WRITE > 0 {
		joins = append(joins, "PROT_WRITE")
	}
	if arg.Raw()&unix.PROT_EXEC > 0 {
		joins = append(joins, "PROT_EXEC")
	}
	if arg.Raw()&unix.PROT_GROWSUP > 0 {
		joins = append(joins, "PROT_GROWSUP")
	}
	if arg.Raw()&unix.PROT_GROWSDOWN > 0 {
		joins = append(joins, "PROT_GROWSDOWN")
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

package annotation

/*
#include <sys/quota.h>
*/
import "C"

var quotactlCmds = map[uintptr]string{
	C.Q_QUOTAON:  "Q_QUOTAON",
	C.Q_QUOTAOFF: "Q_QUOTAOFF",
}

func AnnotateQuotactlCmd(arg Arg, _ int) {
	if cmd, ok := quotactlCmds[arg.Raw()]; ok {
		arg.SetAnnotation(cmd, true)
	}
}

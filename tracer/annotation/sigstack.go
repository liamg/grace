package annotation

/*
#include <linux/signal.h>
*/
import "C"

var sigStackFlags = map[int]string{
	C.SS_AUTODISARM: "SS_AUTODISARM",
}

func AnnotateSigStackFlags(arg Arg, _ int) {
	for flag, name := range sigStackFlags {
		if int(arg.Raw())&flag != 0 {
			arg.SetAnnotation(name, true)
		}
	}
}

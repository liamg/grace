package annotation

import "fmt"

func AnnotateHex(arg Arg, _ int) {
	arg.SetAnnotation(fmt.Sprintf("0x%x", arg.Raw()), true)
}

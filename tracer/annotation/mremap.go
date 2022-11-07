package annotation

import "strings"

// see https://github.com/torvalds/linux/blob/master/include/uapi/linux/mman.h
func AnnotateMRemapFlags(arg Arg, _ int) {
	flags := map[int]string{
		1: "MREMAP_MAYMOVE",
		2: "MREMAP_FIXED",
		4: "MREMAP_DONTUNMAP",
	}
	var joins []string
	for flag, str := range flags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

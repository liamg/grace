package annotation

import "C"

var numaModeFlags = map[int]string{
	0:  "MPOL_DEFAULT",
	1:  "MPOL_PREFERRED",
	2:  "MPOL_BIND",
	3:  "MPOL_INTERLEAVE",
	4:  "MPOL_LOCAL",
	5:  "MPOL_MAX",
	6:  "MPOL_F_STATIC_NODES",
	7:  "MPOL_F_RELATIVE_NODES",
	8:  "MPOL_MF_STRICT",
	9:  "MPOL_MF_MOVE",
	10: "MPOL_MF_MOVE_ALL",
}

func AnnotateNumaModeFlag(arg Arg, pid int) {
	if str, ok := numaModeFlags[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

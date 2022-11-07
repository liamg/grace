package annotation

/*
#include <sys/prctl.h>

#ifndef _ASM_X86_PRCTL_H
#define _ASM_X86_PRCTL_H
#define ARCH_SET_GS			0x1001
#define ARCH_SET_FS			0x1002
#define ARCH_GET_FS			0x1003
#define ARCH_GET_GS			0x1004
#endif

int iARCH_SET_GS = ARCH_SET_GS;
int iARCH_SET_FS = ARCH_SET_FS;
int iARCH_GET_FS = ARCH_GET_FS;
int iARCH_GET_GS = ARCH_GET_GS;
*/
import "C"

var archPrCodes = map[int]string{
	int(C.iARCH_SET_GS): "ARCH_SET_GS",
	int(C.iARCH_SET_FS): "ARCH_SET_FS",
	int(C.iARCH_GET_FS): "ARCH_GET_FS",
	int(C.iARCH_GET_GS): "ARCH_GET_GS",
}

func AnnotateArchPrctrlCode(arg Arg, pid int) {
	if s, ok := archPrCodes[int(arg.Raw())]; ok {
		arg.SetAnnotation(s, true)
	} else {
		AnnotateHex(arg, pid)
	}
}

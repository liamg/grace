package annotation

/*
#include <sys/types.h>
#include <sys/resource.h>
#include <sys/wait.h>

int fWIFEXITED(int status) {
	return WIFEXITED(status);
}
int fWEXITSTATUS(int status) {
	return WEXITSTATUS(status);
}
int fWIFSIGNALED(int status) {
	return WIFSIGNALED(status);
}
int fWTERMSIG(int status) {
	return WTERMSIG(status);
}
int fWIFSTOPPED(int status) {
	return WIFSTOPPED(status);
}
int fWSTOPSIG(int status) {
	return WSTOPSIG(status);
}
*/
import "C"

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

var waitOptions = map[int]string{
	unix.WNOHANG:    "WNOHANG",
	unix.WUNTRACED:  "WUNTRACED",
	unix.WCONTINUED: "WCONTINUED",
	unix.WEXITED:    "WEXITED",
	unix.WNOWAIT:    "WNOWAIT",
	unix.WALL:       "__WALL",
	unix.WNOTHREAD:  "__WNOTHREAD",
	unix.WCLONE:     "__WCLONE",
}

func AnnotateWaitOptions(arg Arg, _ int) {
	var joins []string
	for flag, name := range waitOptions {
		if arg.Raw()&uintptr(flag) != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

func AnnotateWaitStatus(arg Arg, _ int) {
	status := syscall.WaitStatus(arg.Raw())

	var joins []string

	if C.fWIFEXITED(C.int(status)) != 0 {
		joins = append(joins, "WIFEXITED(s)")
		result := C.fWEXITSTATUS(C.int(status))
		joins = append(joins, fmt.Sprintf("WEXITSTATUS(s) == %d", result))
	}

	if C.fWIFSIGNALED(C.int(status)) != 0 {
		joins = append(joins, "WIFSIGNALED(s)")
		result := int(C.fWTERMSIG(C.int(status)))
		joins = append(joins, fmt.Sprintf("WTERMSIG(s) == %s", SignalToString(result)))
	}

	if C.fWIFSTOPPED(C.int(status)) != 0 {
		joins = append(joins, "WIFSTOPPED(s)")
		result := int(C.fWSTOPSIG(C.int(status)))
		joins = append(joins, fmt.Sprintf("WSTOPSIG(s) == %s", SignalToString(result)))
	}

	arg.SetAnnotation(strings.Join(joins, " && "), len(joins) > 0)
}

var waitIDs = map[int]string{
	unix.P_PID:  "P_PID",
	unix.P_PGID: "P_PGID",
	unix.P_ALL:  "P_ALL",
}

func AnnotateIDType(arg Arg, _ int) {
	if str, ok := waitIDs[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

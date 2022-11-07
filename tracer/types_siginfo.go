package tracer

import (
	"syscall"
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"

	"golang.org/x/sys/unix"
)

/*
#include <signal.h>
*/
import "C"

func init() {
	registerTypeHandler(argTypeSigInfo, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {

			raw, err := readSize(pid, raw, unsafe.Sizeof(unix.Siginfo{}))
			if err != nil {
				return err
			}

			var info unix.Siginfo
			if err := decodeStruct(raw, &info); err != nil {
				return err
			}

			arg.obj = convertSigInfo(&info, pid)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertSigInfo(u *unix.Siginfo, pid int) *Object {

	signo := Arg{
		name: "signo",
		t:    ArgTypeInt,
		raw:  uintptr(u.Signo),
	}
	annotation.AnnotateSignal(&signo, pid)

	errStr := annotation.ErrNoToString(int(u.Errno))
	errno := Arg{
		name:       "errno",
		t:          ArgTypeInt,
		raw:        uintptr(u.Errno),
		annotation: errStr,
		replace:    errStr != "",
	}

	codeStr := annotation.SignalCodeToString(syscall.Signal(u.Signo), u.Code)

	code := Arg{
		name:       "code",
		t:          ArgTypeInt,
		raw:        uintptr(u.Code),
		annotation: codeStr,
		replace:    codeStr != "",
	}

	return &Object{
		Name: "siginfo",
		Properties: []Arg{
			signo,
			errno,
			code,
		},
	}
}

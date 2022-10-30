package tracer

import (
	"strings"
	"unsafe"
)

/*
#include <signal.h>

const void* sig_dfl = SIG_DFL;
const void* sig_ign = SIG_IGN;
*/
import "C"

type sigAction struct {
	Handler   uintptr
	Sigaction uintptr
	Mask      int
	Flags     int
	Restorer  uintptr
}

func init() {
	registerTypeHandler(ArgTypeSigAction, func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
		if raw > 0 {

			raw, err := readSize(pid, raw, unsafe.Sizeof(sigAction{}))
			if err != nil {
				return err
			}

			var action sigAction
			if err := decodeStruct(raw, &action); err != nil {
				return err
			}

			arg.obj = convertSigAction(&action)
		}
		return nil
	})
}

func convertSigAction(action *sigAction) *Object {
	obj := Object{
		Name: "sigaction",
	}

	var handlerStr string
	switch action.Handler {
	case uintptr(C.sig_dfl):
		handlerStr = "SIG_DFL"
	case uintptr(C.sig_ign):
		handlerStr = "SIG_IGN"
	}
	obj.Properties = append(obj.Properties, Arg{
		name:       "handler",
		t:          ArgTypeAddress,
		raw:        action.Handler,
		annotation: handlerStr,
		replace:    handlerStr != "",
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "sigaction",
		t:    ArgTypeAddress,
		raw:  action.Sigaction,
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "mask",
		t:    ArgTypeInt,
		raw:  uintptr(action.Mask),
	})

	var signActionFlags []string
	if action.Flags&C.SA_NOCLDSTOP != 0 {
		signActionFlags = append(signActionFlags, "SA_NOCLDSTOP")
	}
	if action.Flags&C.SA_NOCLDWAIT != 0 {
		signActionFlags = append(signActionFlags, "SA_NOCLDWAIT")
	}
	if action.Flags&C.SA_NODEFER != 0 {
		signActionFlags = append(signActionFlags, "SA_NODEFER")
	}
	if action.Flags&C.SA_ONSTACK != 0 {
		signActionFlags = append(signActionFlags, "SA_ONSTACK")
	}
	if action.Flags&C.SA_RESETHAND != 0 {
		signActionFlags = append(signActionFlags, "SA_RESETHAND")
	}
	if action.Flags&C.SA_RESTART != 0 {
		signActionFlags = append(signActionFlags, "SA_RESTART")
	}
	if action.Flags&C.SA_RESTORER != 0 {
		signActionFlags = append(signActionFlags, "SA_RESTORER")
	}
	if action.Flags&C.SA_SIGINFO != 0 {
		signActionFlags = append(signActionFlags, "SA_SIGINFO")
	}
	if action.Flags&C.SA_UNSUPPORTED != 0 {
		signActionFlags = append(signActionFlags, "SA_UNSUPPORTED")
	}
	if action.Flags&C.SA_EXPOSE_TAGBITS != 0 {
		signActionFlags = append(signActionFlags, "SA_EXPOSE_TAGBITS")
	}
	flagStr := strings.Join(signActionFlags, "|")

	obj.Properties = append(obj.Properties, Arg{
		name:       "flags",
		t:          ArgTypeInt,
		raw:        uintptr(action.Flags),
		annotation: flagStr,
		replace:    flagStr != "",
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "restorer",
		t:    ArgTypeAddress,
		raw:  action.Restorer,
	})

	return &obj
}

func annotateSigProcMaskFlags(arg *Arg, pid int) {

	var joins []string
	if arg.Raw()&C.SIG_BLOCK != 0 {
		joins = append(joins, "SIG_BLOCK")
	}
	if arg.Raw()&C.SIG_UNBLOCK != 0 {
		joins = append(joins, "SIG_UNBLOCK")
	}
	if arg.Raw()&C.SIG_SETMASK != 0 {
		joins = append(joins, "SIG_SETMASK")
	}

	arg.annotation = strings.Join(joins, "|")
	arg.replace = arg.annotation != ""
}

package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeItimerspec, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(unix.ItimerSpec{}))
			if err != nil {
				return err
			}

			var timeVal unix.ItimerSpec
			if err := decodeStruct(rawTimeVal, &timeVal); err != nil {
				return err
			}

			arg.obj = convertITimerSpec(&timeVal)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertITimerSpec(u *unix.ItimerSpec) *Object {
	if u == nil {
		return nil
	}
	return &Object{
		Name: "itimerspec",
		Properties: []Arg{
			{
				name: "interval",
				t:    ArgTypeObject,
				obj:  convertTimeSpec(&u.Interval),
			},
			{
				name: "value",
				t:    ArgTypeObject,
				obj:  convertTimeSpec(&u.Value),
			},
		},
	}
}

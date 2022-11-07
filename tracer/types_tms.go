package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeTms, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Tms{}))
			if err != nil {
				return err
			}

			var times unix.Tms
			if err := decodeStruct(rawVal, &times); err != nil {
				return err
			}

			arg.obj = convertTimes(&times)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertTimes(times *unix.Tms) *Object {
	return &Object{
		Name: "times",
		Properties: []Arg{
			{
				name: "utime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(times.Utime),
			},
			{
				name: "stime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(times.Stime),
			},
			{
				name: "cutime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(times.Cutime),
			},
			{
				name: "cstime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(times.Cstime),
			},
		},
	}
}

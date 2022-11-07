package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeItimerval, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Itimerval{}))
			if err != nil {
				return err
			}

			var timeVal unix.Itimerval
			if err := decodeStruct(rawTimeVal, &timeVal); err != nil {
				return err
			}

			arg.obj = convertITimerVal(&timeVal)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertITimerVal(u *unix.Itimerval) *Object {
	return &Object{
		Name: "itimerval",
		Properties: []Arg{
			{
				name: "interval",
				t:    ArgTypeObject,
				obj:  convertTimeVal(&u.Interval),
			},
			{
				name: "value",
				t:    ArgTypeObject,
				obj:  convertTimeVal(&u.Value),
			},
		},
	}
}

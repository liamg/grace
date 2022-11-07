package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeUtimbuf, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Utimbuf{}))
			if err != nil {
				return err
			}

			var timeVal unix.Utimbuf
			if err := decodeStruct(rawTimeVal, &timeVal); err != nil {
				return err
			}

			arg.obj = convertTimeBuf(&timeVal)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertTimeBuf(u *unix.Utimbuf) *Object {
	return &Object{
		Name: "timbuf",
		Properties: []Arg{
			{
				name: "actime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Actime),
			},
			{
				name: "modtime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Modtime),
			},
		},
	}
}

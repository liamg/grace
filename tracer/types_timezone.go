package tracer

import (
	"unsafe"
)

type timezone struct {
	Minuteswest int32
	Dsttime     int32
}

func init() {
	registerTypeHandler(argTypeTimezone, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(timezone{}))
			if err != nil {
				return err
			}

			var tz timezone
			if err := decodeStruct(mem, &tz); err != nil {
				return err
			}

			arg.obj = convertTimezone(&tz)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertTimezone(tz *timezone) *Object {
	return &Object{
		Name: "timezone",
		Properties: []Arg{
			{
				name: "minuteswest",
				t:    ArgTypeInt,
				raw:  uintptr(tz.Minuteswest),
			},
			{
				name: "dsttime",
				t:    ArgTypeInt,
				raw:  uintptr(tz.Dsttime),
			},
		},
	}
}

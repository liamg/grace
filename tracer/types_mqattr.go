package tracer

import (
	"unsafe"
)

type mqattr struct {
	Flags   int32
	Maxmsg  int32
	Msgsize int32
	Curmsgs int32
}

func init() {
	registerTypeHandler(argTypeMqAttr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(mqattr{}))
			if err != nil {
				return err
			}

			var attr mqattr
			if err := decodeStruct(rawTimeVal, &attr); err != nil {
				return err
			}

			arg.obj = convertMqAttr(attr)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertMqAttr(a mqattr) *Object {
	return &Object{
		Name: "mq_attr",
		Properties: []Arg{
			{
				name: "flags",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(a.Flags),
			},
			{
				name: "maxmsg",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(a.Maxmsg),
			},
			{
				name: "msgsize",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(a.Msgsize),
			},
			{
				name: "curmsgs",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(a.Curmsgs),
			},
		},
	}
}

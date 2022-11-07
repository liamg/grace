package tracer

import (
	"unsafe"
)

type iocb struct {
	Data      uint64
	Key       uint32
	Opcode    uint16
	Priority  uint16
	Flags     uint32
	Fd        uint32
	Offset    uint64
	Addr      uint64
	Len       uint32
	Pos       uint64
	Reserved2 uint32
	Reserved3 uint64
}

func init() {
	registerTypeHandler(argTypeIoCB, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {

			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(ioevent{}))
			if err != nil {
				return err
			}

			var rawIoCB iocb
			if err := decodeStruct(mem, &rawIoCB); err != nil {
				return err
			}

			arg.t = ArgTypeObject
			arg.obj = convertIoCB(rawIoCB)
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertIoCB(cb iocb) *Object {
	return &Object{
		Name: "iocb",
		Properties: []Arg{
			{
				name: "data",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(cb.Data),
			},
			{
				name: "key",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Key),
			},
			{
				name: "opcode",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Opcode),
			},
			{
				name: "priority",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Priority),
			},
			{
				name: "flags",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Flags),
			},
			{
				name: "fd",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Fd),
			},
			{
				name: "offset",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(cb.Offset),
			},
			{
				name: "addr",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(cb.Addr),
			},
			{
				name: "len",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Len),
			},
			{
				name: "pos",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(cb.Pos),
			},
			{
				name: "reserved2",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(cb.Reserved2),
			},
			{
				name: "reserved3",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(cb.Reserved3),
			},
		},
	}
}

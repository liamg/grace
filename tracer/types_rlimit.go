package tracer

import (
	"syscall"
	"unsafe"
)

func init() {
	registerTypeHandler(argTypeRLimit, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			data, err := readSize(pid, raw, unsafe.Sizeof(syscall.Rlimit{}))
			if err != nil {
				return err
			}

			var rlimit syscall.Rlimit
			if err := decodeStruct(data, &rlimit); err != nil {
				return err
			}

			arg.obj = convertRLimit(&rlimit)
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertRLimit(rlimit *syscall.Rlimit) *Object {
	return &Object{
		Name: "rlimit",
		Properties: []Arg{
			{
				name: "cur",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(rlimit.Cur),
			},
			{
				name: "max",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(rlimit.Max),
			},
		},
	}
}

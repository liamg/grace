package tracer

import (
	"unsafe"
)

func init() {
	registerTypeHandler(argTypeUnsignedIntPtr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		var underlying uint32
		if buf, err := readSize(pid, next, unsafe.Sizeof(underlying)); err == nil {
			arg.raw = uintptr(decodeInt(buf))
			arg.t = ArgTypeUnsignedInt
		} else {
			arg.t = ArgTypeAddress
			if raw == 0 {
				arg.annotation = "NULL"
				arg.replace = true
			}
		}
		return nil
	})
	registerTypeHandler(argTypeUnsignedInt64Ptr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		var underlying uint64
		if buf, err := readSize(pid, next, unsafe.Sizeof(underlying)); err == nil {
			arg.raw = uintptr(decodeInt(buf))
			arg.t = ArgTypeUnsignedInt
		} else {
			arg.t = ArgTypeAddress
			if raw == 0 {
				arg.annotation = "NULL"
				arg.replace = true
			}
		}
		return nil
	})
}

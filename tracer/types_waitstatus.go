package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeWaitStatus, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			bytes, err := readSize(pid, raw, unsafe.Sizeof(unix.WaitStatus(0)))
			if err != nil {
				return err
			}

			status := uintptr(decodeUint(bytes))

			arg.raw = status
			arg.t = ArgTypeUnsignedInt
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

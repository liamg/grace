package tracer

import (
	"reflect"
	"unsafe"
)

func init() {
	registerTypeHandler(ArgTypeIovecArray, func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
		// read the raw C struct from the process memory
		mem, err := readSize(pid, raw, next*unsafe.Sizeof(iovec{}))
		if err != nil {
			return err
		}

		iovecs := make([]iovec, next)
		if err := decodeAnonymous(reflect.ValueOf(iovecs), mem); err != nil {
			return err
		}

		arg.array = convertIovecs(iovecs)
		return nil
	})
}

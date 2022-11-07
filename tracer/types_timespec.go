package tracer

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeTimespec, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Timespec{}))
			if err != nil {
				return err
			}

			var timeVal unix.Timespec
			if err := decodeStruct(rawTimeVal, &timeVal); err != nil {
				return err
			}

			arg.obj = convertTimeSpec(&timeVal)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
	registerTypeHandler(argTypeTimespecArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {

			var size int

			switch metadata.LenSource {
			case LenSourceFixed:
				size = metadata.FixedCount
			default:
				return fmt.Errorf("unsupported count_from: %d", metadata.LenSource)
			}

			rawVals, err := readSize(pid, raw, uintptr(size)*unsafe.Sizeof(unix.Timespec{}))
			if err != nil {
				return err
			}

			vals := make([]unix.Timespec, size)
			if err := decodeAnonymous(reflect.ValueOf(&vals).Elem(), rawVals); err != nil {
				return err
			}

			for _, val := range vals {
				arg.array = append(arg.array, Arg{
					t:   ArgTypeObject,
					obj: convertTimeSpec(&val),
				})
			}
		}

		arg.t = ArgTypeArray
		return nil
	})
}

func convertTimeSpec(u *unix.Timespec) *Object {
	if u == nil {
		return nil
	}
	return &Object{
		Name: "timespec",
		Properties: []Arg{
			{
				name: "sec",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Sec),
			},
			{
				name: "usec",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Nsec),
			},
		},
	}
}

package tracer

import (
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"
)

// used by clone3 syscall
type openHow struct {
	Flags   uint64
	Mode    uint64
	Resolve uint64
}

func init() {
	registerTypeHandler(argTypeOpenHow, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(openHow{}))
			if err != nil {
				return err
			}

			var how openHow
			if err := decodeStruct(mem, &how); err != nil {
				return err
			}

			arg.obj = convertOpenHow(how)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertOpenHow(how openHow) *Object {

	flags := Arg{
		name: "flags",
		t:    ArgTypeUnsignedLong,
		raw:  uintptr(how.Flags),
	}
	annotation.AnnotateOpenFlags(&flags, 0)

	mode := Arg{
		name: "flags",
		t:    ArgTypeUnsignedLong,
		raw:  uintptr(how.Mode),
	}
	annotation.AnnotateAccMode(&mode, 0)

	resolve := Arg{
		name: "resolve",
		t:    ArgTypeUnsignedLong,
		raw:  uintptr(how.Resolve),
	}
	annotation.AnnotateResolveFlags(&mode, 0)

	return &Object{
		Name: "open_how",
		Properties: []Arg{
			flags,
			mode,
			resolve,
		},
	}
}

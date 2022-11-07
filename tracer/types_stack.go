package tracer

import (
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"
)

type stack struct {
	Sp    uintptr
	Flags uint32
	Size  uint32
}

func init() {
	registerTypeHandler(argTypeStack, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(stack{}))
			if err != nil {
				return err
			}

			var stack stack
			if err := decodeStruct(mem, &stack); err != nil {
				return err
			}

			arg.obj = convertStack(&stack)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertStack(st *stack) *Object {

	flagsArg := Arg{
		name: "flags",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(st.Flags),
	}
	annotation.AnnotateSigStackFlags(&flagsArg, 0)

	return &Object{
		Name: "stack",
		Properties: []Arg{
			{
				name: "sp",
				t:    ArgTypeAddress,
				raw:  st.Sp,
			},
			flagsArg,
			{
				name: "size",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(st.Size),
			},
		},
	}
}

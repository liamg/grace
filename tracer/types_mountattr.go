package tracer

import (
	"unsafe"
)

// used by clone3 syscall
type mountAttr struct {
	Set             uint64
	Clear           uint64
	Propagation     uint64
	UserNamespaceFd uint64
}

func init() {
	registerTypeHandler(argTypeMountAttr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(mountAttr{}))
			if err != nil {
				return err
			}

			var attr mountAttr
			if err := decodeStruct(mem, &attr); err != nil {
				return err
			}

			arg.obj = convertMountAttr(attr)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertMountAttr(attr mountAttr) *Object {
	return &Object{
		Name: "mount_attr",
		Properties: []Arg{
			{
				name: "set",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(attr.Set),
			},
			{
				name: "clear",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(attr.Clear),
			},
			{
				name: "propagation",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(attr.Propagation),
			},
			{
				name: "user_namespace_fd",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(attr.UserNamespaceFd),
			},
		},
	}
}

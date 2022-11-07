package tracer

import (
	"unsafe"
)

// used by clone3 syscall
type landlockRulesetAttr struct {
	HandledAccessFS uint64
}

func init() {
	registerTypeHandler(argTypeLandlockRulesetAttr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(landlockRulesetAttr{}))
			if err != nil {
				return err
			}

			var attr landlockRulesetAttr
			if err := decodeStruct(mem, &attr); err != nil {
				return err
			}

			arg.obj = convertLandlockRulesetAttr(attr)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertLandlockRulesetAttr(attr landlockRulesetAttr) *Object {
	return &Object{
		Name: "landlock_ruleset_attr",
		Properties: []Arg{
			{
				name: "handled_access_fs",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(attr.HandledAccessFS),
			},
		},
	}
}

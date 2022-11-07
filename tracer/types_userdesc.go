package tracer

import (
	"unsafe"
)

type ldt struct {
	EntryNumber uint32
	BaseAddr    uint32
	Limit       uint32
}

func init() {
	registerTypeHandler(argTypeUserDesc, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(ldt{}))
			if err != nil {
				return err
			}

			var table ldt
			if err := decodeStruct(rawVal, &table); err != nil {
				return err
			}

			arg.obj = convertLDT(table)
			if err != nil {
				return err
			}
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertLDT(table ldt) *Object {
	return &Object{
		Name: "ldt",
		Properties: []Arg{
			{
				name: "entry_number",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(table.EntryNumber),
			},
			{
				name: "base_addr",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(table.BaseAddr),
			},
			{
				name: "limit",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(table.Limit),
			},
		},
	}
}

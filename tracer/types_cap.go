package tracer

import (
	"unsafe"
)

type capHeader struct {
	Version uint32
	Pid     int32
}

type capData struct {
	Effective   uint32
	Permitted   uint32
	Inheritable uint32
}

func init() {
	registerTypeHandler(argTypeCapUserHeader, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			mem, err := readSize(pid, raw, unsafe.Sizeof(capHeader{}))
			if err != nil {
				return err
			}
			var cap capHeader
			if err := decodeStruct(mem, &cap); err != nil {
				return err
			}
			arg.obj = convertCapHeader(&cap)
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
	registerTypeHandler(argTypeCapUserData, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			mem, err := readSize(pid, raw, unsafe.Sizeof(capData{}))
			if err != nil {
				return err
			}
			var cap capData
			if err := decodeStruct(mem, &cap); err != nil {
				return err
			}
			arg.obj = convertCapData(&cap)
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertCapHeader(cap *capHeader) *Object {
	return &Object{
		Name: "hdr",
		Properties: []Arg{
			{
				name: "version",
				t:    ArgTypeInt,
				raw:  uintptr(cap.Version),
			},
			{
				name: "pid",
				t:    ArgTypeInt,
				raw:  uintptr(cap.Pid),
			},
		},
	}
}

func convertCapData(cap *capData) *Object {
	return &Object{
		Name: "data",
		Properties: []Arg{
			{
				name: "effective",
				t:    ArgTypeInt,
				raw:  uintptr(cap.Effective),
			},
			{
				name: "permitted",
				t:    ArgTypeInt,
				raw:  uintptr(cap.Permitted),
			},
			{
				name: "inheritable",
				t:    ArgTypeInt,
				raw:  uintptr(cap.Inheritable),
			},
		},
	}
}

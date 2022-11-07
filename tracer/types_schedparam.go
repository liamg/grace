package tracer

import (
	"unsafe"
)

type schedParam struct {
	Priority int32
}

func init() {
	registerTypeHandler(argTypeSchedParam, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(schedParam{}))
			if err != nil {
				return err
			}

			var param schedParam
			if err := decodeStruct(rawVal, &param); err != nil {
				return err
			}

			arg.obj = convertSchedParam(param)
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

func convertSchedParam(param schedParam) *Object {
	return &Object{
		Name: "sched_param",
		Properties: []Arg{
			{
				name: "priority",
				t:    ArgTypeInt,
				raw:  uintptr(param.Priority),
			},
		},
	}
}

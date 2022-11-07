package tracer

import (
	"unsafe"
)

type schedAttr struct {
	Size          uint32
	SchedPolicy   uint32
	SchedFlags    uint64
	SchedNice     int32
	SchedPriority uint32
	SchedRuntime  uint64
	SchedDeadline uint64
	SchedPeriod   uint64
}

func init() {
	registerTypeHandler(argTypeSchedAttr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(schedAttr{}))
			if err != nil {
				return err
			}

			var attr schedAttr
			if err := decodeStruct(rawVal, &attr); err != nil {
				return err
			}

			arg.obj = convertSchedAttr(attr)
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

func convertSchedAttr(param schedAttr) *Object {
	return &Object{
		Name: "sched_attr",
		Properties: []Arg{
			{
				name: "size",
				t:    ArgTypeInt,
				raw:  uintptr(param.Size),
			},
			{
				name: "sched_policy",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedPolicy),
			},
			{
				name: "sched_flags",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedFlags),
			},
			{
				name: "sched_nice",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedNice),
			},
			{
				name: "sched_priority",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedPriority),
			},
			{
				name: "sched_runtime",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedRuntime),
			},
			{
				name: "sched_deadline",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedDeadline),
			},
			{
				name: "sched_period",
				t:    ArgTypeInt,
				raw:  uintptr(param.SchedPeriod),
			},
		},
	}
}

package tracer

import (
	"unsafe"
)

type iouringParams struct {
	SqEntries    uint32
	CqEntries    uint32
	Flags        uint32
	SqThreadCpu  uint32
	SqThreadIdle uint32
	Features     uint32
	WqFd         uint32
}

func init() {
	registerTypeHandler(argTypeIoUringParams, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(iouringParams{}))
			if err != nil {
				return err
			}

			var params iouringParams
			if err := decodeStruct(rawTimeVal, &params); err != nil {
				return err
			}

			arg.obj = convertIoUringParams(params)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertIoUringParams(u iouringParams) *Object {
	return &Object{
		Name: "io_uring_params",
		Properties: []Arg{
			{
				name: "sq_entries",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.SqEntries),
			},
			{
				name: "cq_entries",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.CqEntries),
			},
			{
				name: "flags",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Flags),
			},
			{
				name: "sq_thread_cpu",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.SqThreadCpu),
			},
			{
				name: "sq_thread_idle",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.SqThreadIdle),
			},
			{
				name: "features",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.Features),
			},
			{
				name: "wq_fd",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(u.WqFd),
			},
		},
	}
}

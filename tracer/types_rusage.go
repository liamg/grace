package tracer

import (
	"syscall"
	"unsafe"
)

func init() {
	registerTypeHandler(argTypeRUsage, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			data, err := readSize(pid, raw, unsafe.Sizeof(syscall.Rusage{}))
			if err != nil {
				return err
			}

			var rusage syscall.Rusage
			if err := decodeStruct(data, &rusage); err != nil {
				return err
			}

			arg.obj = convertRUsage(&rusage)
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertRUsage(rusage *syscall.Rusage) *Object {
	return &Object{
		Name: "rusage",
		Properties: []Arg{
			{
				name: "utime",
				t:    ArgTypeObject,
				obj: &Object{
					Name: "timeval",
					Properties: []Arg{
						{
							name: "tv_sec",
							t:    ArgTypeLong,
							raw:  uintptr(rusage.Utime.Sec),
						},
						{
							name: "tv_usec",
							t:    ArgTypeLong,
							raw:  uintptr(rusage.Utime.Usec),
						},
					},
				},
			},
			{
				name: "stime",
				t:    ArgTypeObject,
				obj: &Object{
					Name: "timeval",
					Properties: []Arg{
						{
							name: "tv_sec",
							t:    ArgTypeLong,
							raw:  uintptr(rusage.Stime.Sec),
						},
						{
							name: "tv_usec",
							t:    ArgTypeLong,
							raw:  uintptr(rusage.Stime.Usec),
						},
					},
				},
			},
			{
				name: "maxrss",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Maxrss),
			},
			{
				name: "ixrss",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Ixrss),
			},
			{
				name: "idrss",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Idrss),
			},
			{
				name: "isrss",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Isrss),
			},
			{
				name: "minflt",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Minflt),
			},
			{
				name: "majflt",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Majflt),
			},
			{
				name: "nswap",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Nswap),
			},
			{
				name: "inblock",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Inblock),
			},
			{
				name: "oublock",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Oublock),
			},
			{
				name: "msgsnd",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Msgsnd),
			},
			{
				name: "msgrcv",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Msgrcv),
			},
			{
				name: "nsignals",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Nsignals),
			},
			{
				name: "nvcsw",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Nvcsw),
			},
			{
				name: "nivcsw",
				t:    ArgTypeLong,
				raw:  uintptr(rusage.Nivcsw),
			},
		},
	}
}

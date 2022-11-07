package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeSysinfo, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Sysinfo_t{}))
			if err != nil {
				return err
			}

			var sysinfo unix.Sysinfo_t
			if err := decodeStruct(rawVal, &sysinfo); err != nil {
				return err
			}

			arg.obj, err = convertSysinfo(&sysinfo, pid)
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

func convertLoadsToArray(loads [3]uint64) []Arg {
	var loadsArray []Arg
	for _, load := range loads {
		loadsArray = append(loadsArray, Arg{
			name: "load",
			t:    ArgTypeUnsignedLong,
			raw:  uintptr(load),
		})
	}
	return loadsArray
}

func convertSysinfo(sysinfo *unix.Sysinfo_t, _ int) (*Object, error) {
	return &Object{
		Name: "sysinfo",
		Properties: []Arg{
			{
				name: "uptime",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Uptime),
			},
			{
				name:  "loads",
				t:     ArgTypeArray,
				array: convertLoadsToArray(sysinfo.Loads),
			},
			{
				name: "totalram",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Totalram),
			},
			{
				name: "freeram",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Freeram),
			},
			{
				name: "sharedram",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Sharedram),
			},
			{
				name: "bufferram",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Bufferram),
			},
			{
				name: "totalswap",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Totalswap),
			},
			{
				name: "freeswap",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Freeswap),
			},
			{
				name: "procs",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(sysinfo.Procs),
			},
			{
				name: "totalhigh",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Totalhigh),
			},
			{
				name: "freehigh",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(sysinfo.Freehigh),
			},
			{
				name: "mem_unit",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(sysinfo.Unit),
			},
		},
	}, nil
}

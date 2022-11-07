package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeStatfs, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Statfs_t{}))
			if err != nil {
				return err
			}

			var stat unix.Statfs_t
			if err := decodeStruct(rawVal, &stat); err != nil {
				return err
			}

			arg.obj, err = convertStatfs(&stat, pid)
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

func convertStatfs(stat *unix.Statfs_t, _ int) (*Object, error) {
	return &Object{
		Name: "statfs",
		Properties: []Arg{
			{
				name: "type",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Type),
			},
			{
				name: "bsize",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Bsize),
			},

			{
				name: "blocks",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Blocks),
			},
			{
				name: "bfree",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Bfree),
			},
			{
				name: "bavail",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Bavail),
			},
			{
				name: "files",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Files),
			},
			{
				name: "ffree",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Ffree),
			},
			{
				name: "namelen",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Namelen),
			},
			{
				name: "frsize",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Frsize),
			},
			{
				name: "flags",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Flags),
			},
		},
	}, nil
}

package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeStatX, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			// read the raw C struct from the process memory
			rawStat, err := readSize(pid, raw, unsafe.Sizeof(unix.Statx_t{}))
			if err != nil {
				return err
			}

			var stat unix.Statx_t
			if err := decodeStruct(rawStat, &stat); err != nil {
				return err
			}

			arg.obj = convertStatx(&stat)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertStatx(stat *unix.Statx_t) *Object {
	obj := Object{
		Name: "statx",
	}

	obj.Properties = append(obj.Properties, Arg{
		name: "mask",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Mask),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "blksize",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Blksize),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "attributes",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Attributes),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "nlink",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Nlink),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "uid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Uid),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "gid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Gid),
	})

	obj.Properties = append(obj.Properties, Arg{
		name:       "mode",
		t:          ArgTypeUnsignedInt,
		raw:        uintptr(stat.Mode),
		annotation: permModeToString(uint32(stat.Mode)),
		replace:    true,
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "ino",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Ino),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "size",
		t:    ArgTypeInt,
		raw:  uintptr(stat.Size),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "blocks",
		t:    ArgTypeInt,
		raw:  uintptr(stat.Blocks),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "attributes_mask",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Attributes_mask),
	})

	// -

	obj.Properties = append(obj.Properties, Arg{
		name: "atime",
		t:    ArgTypeObject,
		obj:  convertStatXTimestamp(stat.Atime),
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "btime",
		t:    ArgTypeObject,
		obj:  convertStatXTimestamp(stat.Btime),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "ctime",
		t:    ArgTypeObject,
		obj:  convertStatXTimestamp(stat.Ctime),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "mtime",
		t:    ArgTypeObject,
		obj:  convertStatXTimestamp(stat.Mtime),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "rdev_major",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Rdev_major),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "rdev_minor",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Rdev_minor),
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "dev_major",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Dev_major),
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "dev_minor",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Dev_minor),
	})
	obj.Properties = append(obj.Properties, Arg{
		name: "mnt_id",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Mnt_id),
	})

	return &obj
}

func convertStatXTimestamp(t unix.StatxTimestamp) *Object {
	return &Object{
		Name: "statx_timestamp",
		Properties: []Arg{
			{
				name: "sec",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(t.Sec),
			},
			{
				name: "nsec",
				t:    ArgTypeUnsignedInt,
				raw:  uintptr(t.Nsec),
			},
		},
	}
}

package tracer

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeStat, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			// read the raw C struct from the process memory
			rawStat, err := readSize(pid, raw, unsafe.Sizeof(syscall.Stat_t{}))
			if err != nil {
				return err
			}

			var stat syscall.Stat_t
			if err := decodeStruct(rawStat, &stat); err != nil {
				return err
			}

			arg.obj = convertStat(&stat)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertStat(stat *syscall.Stat_t) *Object {
	obj := Object{
		Name: "stat",
	}

	obj.Properties = append(obj.Properties, Arg{
		name:       "mode",
		t:          ArgTypeUnsignedInt,
		raw:        uintptr(stat.Mode),
		annotation: permModeToString(stat.Mode),
		replace:    true,
	})

	obj.Properties = append(obj.Properties, Arg{
		name:       "dev",
		t:          ArgTypeUnsignedInt,
		raw:        uintptr(stat.Dev),
		annotation: annotation.DeviceToString(stat.Dev),
		replace:    true,
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "ino",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Ino),
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
		name: "blksize",
		t:    ArgTypeInt,
		raw:  uintptr(stat.Blksize),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "blocks",
		t:    ArgTypeInt,
		raw:  uintptr(stat.Blocks),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "size",
		t:    ArgTypeInt,
		raw:  uintptr(stat.Size),
	})

	obj.Properties = append(obj.Properties, Arg{
		name: "nlink",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(stat.Nlink),
	})

	obj.Properties = append(obj.Properties, Arg{
		name:       "rdev",
		t:          ArgTypeUnsignedInt,
		raw:        uintptr(stat.Rdev),
		annotation: annotation.DeviceToString(stat.Rdev),
		replace:    true,
	})

	return &obj
}

func permModeToString(mode uint32) string {

	flags := map[uint32]string{
		unix.S_IFBLK:  "S_IFBLK",
		unix.S_IFCHR:  "S_IFCHR",
		unix.S_IFIFO:  "S_IFIFO",
		unix.S_IFLNK:  "S_IFLNK",
		unix.S_IFREG:  "S_IFREG",
		unix.S_IFDIR:  "S_IFDIR",
		unix.S_IFSOCK: "S_IFSOCK",
		unix.S_ISUID:  "S_ISUID",
		unix.S_ISGID:  "S_ISGID",
		unix.S_ISVTX:  "S_ISVTX",
	}

	var joins []string
	for flag, name := range flags {
		if mode&syscall.S_IFMT == flag {
			joins = append(joins, name)
		}
	}

	perm := fmt.Sprintf("%04o", int(mode)&0777)

	joins = append(joins, perm)

	return strings.Join(joins, "|")
}

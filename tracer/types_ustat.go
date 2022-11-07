package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeUstat, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(unix.Ustat_t{}))
			if err != nil {
				return err
			}

			var stat unix.Ustat_t
			if err := decodeStruct(rawVal, &stat); err != nil {
				return err
			}

			arg.obj, err = convertUstat(&stat, pid)
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

func convertUstat(stat *unix.Ustat_t, _ int) (*Object, error) {

	var fname []byte
	for _, b := range stat.Fname {
		if b == 0 {
			break
		}
		fname = append(fname, byte(b))
	}

	var fpack []byte
	for _, b := range stat.Fpack {
		if b == 0 {
			break
		}
		fpack = append(fpack, byte(b))
	}

	return &Object{
		Name: "ustat",
		Properties: []Arg{
			{
				name: "tfree",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Tfree),
			},
			{
				name: "tinode",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(stat.Tinode),
			},
			{
				name: "fname",
				t:    ArgTypeData,
				data: fname,
			},
			{
				name: "fpack",
				t:    ArgTypeData,
				data: fpack,
			},
		},
	}, nil
}

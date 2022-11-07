package tracer

import (
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeFdSet, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			// read the raw C struct from the process memory
			rawFdset, err := readSize(pid, raw, unsafe.Sizeof(unix.FdSet{}))
			if err != nil {
				return err
			}

			// safely squish it into a struct in our own memory space
			var fdset unix.FdSet
			if err := decodeStruct(rawFdset, &fdset); err != nil {
				return err
			}

			// convert into a nice array for output
			arg.array = convertFdset(&fdset, pid)
		}
		arg.t = ArgTypeArray
		return nil
	})
}

func convertFdset(fdset *unix.FdSet, pid int) []Arg {
	var fds []Arg
	for _, fd := range fdset.Bits {
		if fd == 0 {
			break
		}
		item := Arg{
			name: "fd",
			t:    ArgTypeInt,
			raw:  uintptr(fd),
		}
		annotation.AnnotateFd(&item, pid)
		fds = append(fds, item)
	}
	return fds
}

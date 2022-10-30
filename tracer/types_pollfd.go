package tracer

import (
	"reflect"
	"unsafe"
)

func init() {
	registerTypeHandler(ArgTypePollFdArray, func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
		if raw > 0 {
			// read the raw C struct from the process memory
			rawPollFds, err := readSize(pid, raw, next*unsafe.Sizeof(pollfd{}))
			if err != nil {
				return err
			}

			pollFds := make([]pollfd, next)
			if err := decodeAnonymous(reflect.ValueOf(pollFds), rawPollFds); err != nil {
				return err
			}

			arg.array = convertPollFds(pollFds)
		}
		return nil
	})
}

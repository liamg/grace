package tracer

import (
	"reflect"
	"unsafe"
)

func init() {
	registerTypeHandler(argTypePollFdArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			// read the raw C struct from the process memory
			rawPollFds, err := readSize(pid, raw, next*unsafe.Sizeof(pollfd{}))
			if err != nil {
				return err
			}

			pollFds := make([]pollfd, next)
			if err := decodeAnonymous(reflect.ValueOf(&pollFds).Elem(), rawPollFds); err != nil {
				return err
			}

			arg.array = convertPollFds(pollFds)
		}
		arg.t = ArgTypeArray
		return nil
	})
}

type pollfd struct {
	Fd      int    /* file descriptor */
	Events  uint16 /* events requested for polling */
	REvents uint32 /* events that occurred during polling */
}

func convertPollFds(fds []pollfd) []Arg {
	var output []Arg
	for _, fd := range fds {
		output = append(output, convertPollFd(fd))
	}
	return output
}

func convertPollFd(fd pollfd) Arg {
	return Arg{
		t: ArgTypeObject,
		obj: &Object{
			Name: "pollfd",
			Properties: []Arg{
				{
					name: "fd",
					t:    ArgTypeInt,
					raw:  uintptr(fd.Fd),
				},
				{
					name: "events",
					t:    ArgTypeInt,
					raw:  uintptr(fd.Events),
				},
				{
					name: "revents",
					t:    ArgTypeInt,
					raw:  uintptr(fd.REvents),
				},
			},
		},
		known: true,
	}
}

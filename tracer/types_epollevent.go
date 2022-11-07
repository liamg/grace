package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeEpollEvent, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(unix.EpollEvent{}))
			if err != nil {
				return err
			}

			var event unix.EpollEvent
			if err := decodeStruct(rawTimeVal, &event); err != nil {
				return err
			}

			arg.obj = convertEpollEvent(&event)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertEpollEvent(ev *unix.EpollEvent) *Object {
	if ev == nil {
		return nil
	}
	return &Object{
		Name: "epoll_event",
		Properties: []Arg{
			{
				name: "events",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(ev.Events),
			},
			{
				name: "fd",
				t:    ArgTypeAddress,
				raw:  uintptr(ev.Fd),
			},
		},
	}
}

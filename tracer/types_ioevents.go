package tracer

import (
	"reflect"
	"unsafe"
)

type ioevent struct {
	Data uint64
	Obj  uint64
	Res  int64
	Res2 int64
}

func init() {
	registerTypeHandler(argTypeIoEvents, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {

			var events []Arg

			// ret contains number of events read
			if ret > 0 {

				// read the raw C struct from the process memory
				mem, err := readSize(pid, raw, unsafe.Sizeof(ioevent{})*ret)
				if err != nil {
					return err
				}

				var rawEvents []ioevent
				if err := decodeAnonymous(reflect.ValueOf(&rawEvents).Elem(), mem); err != nil {
					return err
				}

				events = convertIoEvents(rawEvents)
			}

			arg.t = ArgTypeArray
			arg.array = events
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
	registerTypeHandler(argTypeIoEvent, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {

			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(ioevent{}))
			if err != nil {
				return err
			}

			var rawEvent ioevent
			if err := decodeStruct(mem, &rawEvent); err != nil {
				return err
			}

			arg.t = ArgTypeObject
			arg.obj = convertIoEvent(rawEvent)
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertIoEvents(events []ioevent) []Arg {
	var output []Arg
	for _, event := range events {
		output = append(output, Arg{
			t:   ArgTypeObject,
			obj: convertIoEvent(event),
		})
	}
	return output
}

func convertIoEvent(event ioevent) *Object {
	return &Object{
		Name: "io_event",
		Properties: []Arg{
			{
				name: "data",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(event.Data),
			},
			{
				name: "obj",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(event.Obj),
			},
			{
				name: "res",
				t:    ArgTypeLong,
				raw:  uintptr(event.Res),
			},
			{
				name: "res2",
				t:    ArgTypeLong,
				raw:  uintptr(event.Res2),
			},
		},
	}
}

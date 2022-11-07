package tracer

import (
	"reflect"
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"
)

type sembuf struct {
	Num uint16 /* semaphore index in array */
	Op  int16  /* semaphore operation */
	Flg int16  /* flags for operation */
}

func init() {
	registerTypeHandler(argTypeSembuf, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		// read the raw C struct from the process memory
		mem, err := readSize(pid, raw, next*unsafe.Sizeof(sembuf{})*next)
		if err != nil {
			return err
		}

		sembufs := make([]sembuf, next)
		if err := decodeAnonymous(reflect.ValueOf(&sembufs).Elem(), mem); err != nil {
			return err
		}

		arg.array, err = convertSembufs(sembufs, pid)
		if err != nil {
			return err
		}
		arg.t = ArgTypeArray
		return nil
	})
}

func convertSembufs(bufs []sembuf, pid int) ([]Arg, error) {
	var items []Arg
	for _, buf := range bufs {
		obj := convertSembuf(buf, pid)
		items = append(items, obj)
	}
	return items, nil
}

func convertSembuf(buf sembuf, pid int) Arg {

	flagsArg := Arg{
		name: "flags",
		t:    ArgTypeInt,
		raw:  uintptr(buf.Flg),
	}
	annotation.AnnotateSemFlags(&flagsArg, pid)

	return Arg{
		name: "sops",
		t:    ArgTypeObject,
		obj: &Object{
			Name: "sembuf",
			Properties: []Arg{
				{
					name: "sem_num",
					t:    ArgTypeInt,
					raw:  uintptr(buf.Num),
				},
				{
					name: "sem_op",
					t:    ArgTypeInt,
					raw:  uintptr(buf.Op),
				},
				flagsArg,
			},
		},
	}

}

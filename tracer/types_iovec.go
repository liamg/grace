package tracer

import (
	"reflect"
	"unsafe"
)

func init() {
	registerTypeHandler(argTypeIovecArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		// read the raw C struct from the process memory
		mem, err := readSize(pid, raw, next*unsafe.Sizeof(iovec{}))
		if err != nil {
			return err
		}

		iovecs := make([]iovec, next)
		if err := decodeAnonymous(reflect.ValueOf(&iovecs).Elem(), mem); err != nil {
			return err
		}

		arg.array, err = convertIovecs(iovecs, pid)
		if err != nil {
			return err
		}
		arg.t = ArgTypeArray
		return nil
	})
}

type iovec struct {
	Base uintptr /* Starting address */
	Len  uintptr /* Number of bytes to transfer */
}

func convertIovecs(vecs []iovec, pid int) ([]Arg, error) {
	var output []Arg
	for _, fd := range vecs {
		vec, err := convertIovec(fd, pid)
		if err != nil {
			return nil, err
		}
		output = append(output, *vec)
	}
	return output, nil
}

func convertIovec(vec iovec, pid int) (*Arg, error) {

	base, err := readSize(pid, vec.Base, uintptr(vec.Len))
	if err != nil {
		return nil, err
	}

	return &Arg{
		t: ArgTypeObject,
		obj: &Object{
			Name: "iovec",
			Properties: []Arg{
				{
					name: "base",
					t:    ArgTypeData,
					data: base,
					raw:  vec.Base,
				},
				{
					name: "len",
					t:    ArgTypeUnsignedInt,
					raw:  vec.Len,
				},
			},
		},
		known: true,
	}, nil
}

package tracer

import (
	"fmt"
	"reflect"
)

func init() {
	registerTypeHandler(ArgTypeIntArray, func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
		var count int
		switch metadata.CountFrom {
		case CountLocationNext:
			count = int(next)
		case CountLocationResult:
			count = int(ret)
		case CountLocationFixed:
			count = metadata.FixedCount
		default:
			return fmt.Errorf("syscall %s has no supported count location", metadata.Name)
		}

		mem, err := readSize(pid, raw, 4*uintptr(count))
		if err != nil {
			return err
		}

		target := make([]int32, count)
		if err := decodeAnonymous(reflect.ValueOf(target), mem); err != nil {
			return err
		}

		arg.array = nil
		for i := 0; i < count; i++ {
			arg.array = append(arg.array, Arg{
				t:   ArgTypeInt,
				raw: uintptr(target[i]),
			})
		}

		return nil
	})
}

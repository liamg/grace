package tracer

import (
	"fmt"
	"reflect"
)

func init() {
	registerTypeHandler(argTypeIntArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		var count int
		switch metadata.LenSource {
		case LenSourcePrev:
			count = int(prev)
		case LenSourceNext:
			count = int(next)
		case LenSourceReturnValue:
			count = int(ret)
		case LenSourceFixed:
			count = metadata.FixedCount
		default:
			return fmt.Errorf("syscall %s has no supported count location", metadata.Name)
		}

		mem, err := readSize(pid, raw, 4*uintptr(count))
		if err != nil {
			return err
		}

		target := make([]int32, count)
		if err := decodeAnonymous(reflect.ValueOf(&target).Elem(), mem); err != nil {
			return err
		}

		arg.array = nil
		for i := 0; i < count; i++ {
			arg.array = append(arg.array, Arg{
				t:   ArgTypeInt,
				raw: uintptr(target[i]),
			})
		}
		arg.t = ArgTypeArray

		return nil
	})
}

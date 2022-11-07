package tracer

import (
	"fmt"
	"syscall"
)

func init() {
	registerTypeHandler(ArgTypeData, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		switch metadata.LenSource {
		case LenSourceNextPointer:
			if next == 0 {
				return nil
			}
			if buf, err := readSize(pid, next, 4); err == nil {
				size := uintptr(decodeInt(buf))
				data, err := readSize(pid, raw, size)
				if err != nil {
					return err
				}
				arg.data = data
			}
		case LenSourcePrev:
			data, err := readSize(pid, raw, prev)
			if err != nil {
				return err
			}
			arg.data = data
		case LenSourceNext:
			data, err := readSize(pid, raw, next)
			if err != nil {
				return err
			}
			arg.data = data
		case LenSourceReturnValue:
			data, err := readSize(pid, raw, ret)
			if err != nil {
				return err
			}
			arg.data = data
		default:
			return fmt.Errorf("syscall %s has no supported count location", metadata.Name)
		}
		return nil
	})
}

func readSize(pid int, addr uintptr, size uintptr) ([]byte, error) {

	if size == 0 || size>>(bitSize-1) == 1 { // if negative for this arch
		return nil, nil
	}
	data := make([]byte, size)
	count, err := syscall.PtracePeekData(pid, addr, data)
	if err != nil {
		return nil, fmt.Errorf("read of 0x%x (%d) failed: %w", addr, size, err)
	}
	return data[:count], nil
}

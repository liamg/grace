package tracer

import (
	"fmt"
	"syscall"
)

func init() {
	registerTypeHandler(ArgTypeData, func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
		switch metadata.CountFrom {
		case CountLocationNext:
			data, err := readSize(pid, raw, next)
			if err != nil {
				return err
			}
			arg.data = data
		case CountLocationResult:
			data, err := readSize(pid, raw, ret)
			if err != nil {
				return err
			}
			arg.data = data
		case CountLocationNullTerminator:
			str, err := readString(pid, raw)
			if err != nil {
				return err
			}
			arg.data = []byte(str)
		default:
			return fmt.Errorf("syscall %s has no supported count location", metadata.Name)
		}
		return nil
	})
}

func readSize(pid int, addr uintptr, size uintptr) ([]byte, error) {
	if size == 0 {
		return nil, nil
	}
	data := make([]byte, size)
	count, err := syscall.PtracePeekData(pid, addr, data)
	if err != nil {
		return nil, fmt.Errorf("read of 0x%x (%d) failed: %w", addr, size, err)
	}
	return data[:count], nil
}

func readString(pid int, addr uintptr) (string, error) {
	var output string
	if addr == 0 {
		return output, nil
	}
	data := make([]byte, 1)
	for {
		if _, err := syscall.PtracePeekData(pid, addr, data); err != nil {
			return "", fmt.Errorf("read of 0x%x failed: %w", addr, err)
		}
		if data[0] == 0 {
			break
		}
		output += string(data)
		addr++
	}
	return output, nil
}

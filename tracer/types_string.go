package tracer

import (
	"fmt"
	"syscall"
	"unsafe"
)

/*
	NOTE: an array of strings in C (at least a char *argv[]) is an array of pointers to strings, terminated by a NULL pointer.
*/

func init() {
	registerTypeHandler(argTypeStringArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		var items []Arg
		var offset uintptr
		size := unsafe.Sizeof(uintptr(0))

		for {
			mem, err := readSize(pid, raw+offset, size)
			if err != nil {
				return err
			}
			address := uintptr(decodeUint(mem))
			if address == 0 {
				break
			}
			str, err := readString(pid, address)
			if err != nil {
				return err
			}
			items = append(items, Arg{
				t:    ArgTypeData,
				raw:  address,
				data: []byte(str),
			})
			offset += size
		}
		arg.t = ArgTypeArray
		arg.array = items
		return nil
	})
	registerTypeHandler(argTypeString, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		str, err := readString(pid, raw)
		if err != nil {
			return err
		}
		arg.data = []byte(str)
		arg.t = ArgTypeData
		return nil
	})
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

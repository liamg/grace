package tracer

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"
)

func init() {
	registerTypeHandler(argTypeSockaddr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			rawFamily, err := readSize(pid, raw, unsafe.Sizeof(syscall.RawSockaddrInet4{}.Family))
			if err != nil {
				return err
			}

			family := decodeInt(rawFamily)

			rawSockAddr, err := readSize(pid, raw, next)
			if err != nil {
				return err
			}

			arg.obj, err = convertSockAddr(family, rawSockAddr)
			if err != nil {
				return err
			}
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func htons(port uint16) uint16 {
	return port>>8 | port<<8
}

func convertSockAddr(family int64, raw []byte) (*Object, error) {
	switch family {
	case syscall.AF_INET:
		var target syscall.RawSockaddrInet4
		if err := decodeStruct(raw, &target); err != nil {
			return nil, err
		}
		return &Object{
			Name: "sockaddr",
			Properties: []Arg{
				{
					name:       "family",
					annotation: "AF_INET",
					replace:    true,
				},
				{
					name: "port",
					t:    ArgTypeInt,
					raw:  uintptr(htons(target.Port)),
				},
				{
					name: "addr",
					t:    ArgTypeData,
					data: []byte(fmt.Sprintf("%d.%d.%d.%d", target.Addr[0], target.Addr[1], target.Addr[2], target.Addr[3])),
				},
			},
		}, nil
	case syscall.AF_INET6:
		var target syscall.RawSockaddrInet6
		if err := decodeStruct(raw, &target); err != nil {
			return nil, err
		}
		return &Object{
			Name: "sockaddr",
			Properties: []Arg{
				{
					name:       "family",
					annotation: "AF_INET6",
					replace:    true,
				},
				{
					name: "port",
					t:    ArgTypeInt,
					raw:  uintptr(htons(target.Port)),
				},
				{
					name: "addr",
					t:    ArgTypeData,
					data: []byte(fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", target.Addr[0], target.Addr[1], target.Addr[2], target.Addr[3], target.Addr[4], target.Addr[5], target.Addr[6], target.Addr[7])),
				},
				{
					name: "flowinfo",
					t:    ArgTypeInt,
					raw:  uintptr(target.Flowinfo),
				},
				{
					name: "scope_id",
					t:    ArgTypeInt,
					raw:  uintptr(target.Scope_id),
				},
			},
		}, nil
	case syscall.AF_UNIX:
		var target syscall.RawSockaddrUnix

		if len(raw) < int(unsafe.Sizeof(target)) {
			raw = append(raw, make([]byte, int(unsafe.Sizeof(target))-len(raw))...)
		}

		if err := decodeStruct(raw, &target); err != nil {
			return nil, err
		}
		var path string
		for _, b := range target.Path {
			if b == 0 {
				break
			}
			path += string(byte(b))
		}
		return &Object{
			Name: "sockaddr",
			Properties: []Arg{
				{
					name:       "family",
					annotation: "AF_UNIX",
					replace:    true,
				},
				{
					name: "path",
					t:    ArgTypeData,
					data: []byte(path),
				},
			},
		}, nil
	case syscall.AF_NETLINK:
		var target syscall.RawSockaddrNetlink
		if err := decodeStruct(raw, &target); err != nil {
			return nil, err
		}
		return &Object{
			Name: "sockaddr",
			Properties: []Arg{
				{
					name:       "family",
					annotation: "AF_NETLINK",
					replace:    true,
				},
				{
					name: "pid",
					t:    ArgTypeInt,
					raw:  uintptr(target.Pid),
				},
				{
					name: "groups",
					t:    ArgTypeInt,
					raw:  uintptr(target.Groups),
				},
			},
		}, nil
	default:
		var target syscall.RawSockaddr
		if err := decodeStruct(raw, &target); err != nil {
			return nil, err
		}
		data := make([]byte, len(target.Data))
		for i, b := range target.Data {
			data[i] = byte(b)
		}
		familyStr := annotation.SocketFamily(int(target.Family))
		return &Object{
			Name: "sockaddr",
			Properties: []Arg{
				{
					name:       "family",
					annotation: familyStr,
					replace:    familyStr != "",
				},
				{
					name: "data",
					t:    ArgTypeData,
					data: data,
				},
			},
		}, nil
	}
}

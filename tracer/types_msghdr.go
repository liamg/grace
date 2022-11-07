package tracer

import (
	"reflect"
	"syscall"
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"

	"golang.org/x/sys/unix"
)

// we need our own struct here as the unix.Msghdr struct appears to have a hardcoded 32-bit pointer size???
// is this a bug in Go? or is some preprocessing done per-arch in the syscall pkg?
type msghdr struct {
	Name       uintptr
	Namelen    uintptr
	Iov        uintptr
	Iovlen     uintptr
	Control    uintptr
	Controllen uintptr
	Flags      uintptr
}

type mmsghdr struct {
	MsgHdr msghdr
	MsgLen uintptr
}

func init() {
	registerTypeHandler(argTypeMsghdr, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(msghdr{}))
			if err != nil {
				return err
			}

			var msghdr msghdr
			if err := decodeStruct(rawVal, &msghdr); err != nil {
				return err
			}

			arg.obj, err = convertMsgHdr(&msghdr, pid)
			if err != nil {
				return err
			}
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
	registerTypeHandler(argTypeMMsgHdrArray, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		arg.t = ArgTypeArray
		if raw > 0 {
			rawVal, err := readSize(pid, raw, unsafe.Sizeof(mmsghdr{})*ret)
			if err != nil {
				return err
			}
			var mmsghdrs []mmsghdr
			if err := decodeAnonymous(reflect.ValueOf(&mmsghdrs).Elem(), rawVal); err != nil {
				return err
			}

			arg.array = convertMMsghdrs(mmsghdrs, pid)
		}
		return nil
	})
}

func convertMMsghdrs(mmsghdrs []mmsghdr, pid int) []Arg {

	var args []Arg

	for _, hdr := range mmsghdrs {
		obj, err := convertMsgHdr(&hdr.MsgHdr, pid)
		if err != nil {
			continue
		}
		args = append(args, Arg{
			name: "msg_hdr",
			t:    ArgTypeObject,
			obj: &Object{
				Name: "mmsghdr",
				Properties: []Arg{
					{
						name: "msg_hdr",
						t:    ArgTypeObject,
						obj:  obj,
					},
					{
						name: "msg_len",
						t:    ArgTypeInt,
						raw:  hdr.MsgLen,
					},
				},
			},
		})
	}

	return args
}

func convertMsgHdr(hdr *msghdr, pid int) (*Object, error) {

	rawFamily, err := readSize(pid, hdr.Name, unsafe.Sizeof(syscall.RawSockaddrInet4{}.Family))
	if err != nil {
		return nil, err
	}

	family := decodeInt(rawFamily)

	rawSockAddr, err := readSize(pid, hdr.Name, hdr.Namelen)
	if err != nil {
		return nil, err
	}

	obj, err := convertSockAddr(family, rawSockAddr)
	if err != nil {
		return nil, err
	}

	iovecBytes, err := readSize(pid, hdr.Iov, unsafe.Sizeof(iovec{})*hdr.Iovlen)
	if err != nil {
		return nil, err
	}

	iovecs := make([]iovec, hdr.Iovlen)
	if err := decodeAnonymous(reflect.ValueOf(&iovecs).Elem(), iovecBytes); err != nil {
		return nil, err
	}

	vecs, err := convertIovecs(iovecs, pid)
	if err != nil {
		return nil, err
	}

	controlBytes, err := readSize(pid, hdr.Control, hdr.Controllen)
	if err != nil {
		return nil, err
	}

	var control unix.Cmsghdr
	if err := decodeStruct(controlBytes, &control); err != nil {
		return nil, err
	}

	flags := Arg{
		name: "flags",
		t:    ArgTypeInt,
		raw:  hdr.Flags,
	}
	annotation.AnnotateMsgFlags(&flags, pid)

	return &Object{
		Name: "msghdr",
		Properties: []Arg{
			{
				name: "name",
				t:    ArgTypeObject,
				obj:  obj,
			},
			{
				name:  "iovec",
				t:     ArgTypeArray,
				array: vecs,
			},
			{
				name: "control",
				t:    ArgTypeObject,
				obj:  convertCmsghdr(control, pid),
			},
			flags,
		},
	}, nil
}

func convertCmsghdr(control unix.Cmsghdr, pid int) *Object {

	level := Arg{
		name: "level",
		t:    ArgTypeInt,
		raw:  uintptr(control.Level),
	}

	annotation.AnnotateSocketLevel(&level, pid)

	typ := Arg{
		name: "type",
		t:    ArgTypeInt,
		raw:  uintptr(control.Type),
	}

	annotation.AnnotateControlMessageType(&typ, pid)

	return &Object{
		Name: "cmsghdr",
		Properties: []Arg{
			{
				name: "len",
				t:    ArgTypeInt,
				raw:  uintptr(control.Len),
			},
			level,
			typ,
		},
	}
}

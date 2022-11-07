package tracer

import (
	"syscall"
	"unsafe"
)

func init() {
	registerTypeHandler(argTypeUname, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			data, err := readSize(pid, raw, unsafe.Sizeof(syscall.Utsname{}))
			if err != nil {
				return err
			}

			var uname syscall.Utsname
			if err := decodeStruct(data, &uname); err != nil {
				return err
			}

			arg.obj = convertUname(&uname)
			arg.t = ArgTypeObject
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func getUtsField(f [65]int8) []byte {
	var str []byte
	for i := 0; i < len(f); i++ {
		if f[i] == 0 {
			break
		}
		str = append(str, byte(f[i]))
	}
	return str
}

func convertUname(uname *syscall.Utsname) *Object {
	return &Object{
		Name: "uname",
		Properties: []Arg{
			{
				name: "sysname",
				t:    ArgTypeData,
				data: getUtsField(uname.Sysname),
			},
			{
				name: "nodename",
				t:    ArgTypeData,
				data: getUtsField(uname.Nodename),
			},
			{
				name: "release",
				t:    ArgTypeData,
				data: getUtsField(uname.Release),
			},
			{
				name: "version",
				t:    ArgTypeData,
				data: getUtsField(uname.Version),
			},
			{
				name: "machine",
				t:    ArgTypeData,
				data: getUtsField(uname.Machine),
			},
			{
				name: "domainname",
				t:    ArgTypeData,
				data: getUtsField(uname.Domainname),
			},
		},
	}
}

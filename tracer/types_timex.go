package tracer

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	registerTypeHandler(argTypeTimex, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(unix.Timex{}))
			if err != nil {
				return err
			}

			var tx unix.Timex
			if err := decodeStruct(mem, &tx); err != nil {
				return err
			}

			arg.obj = convertTimex(&tx)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertTimex(tx *unix.Timex) *Object {
	return &Object{
		Name: "timezone",
		Properties: []Arg{
			{
				name: "modes",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Modes),
			},
			{
				name: "offset",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Offset),
			},
			{
				name: "freq",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Freq),
			},
			{
				name: "maxerror",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Maxerror),
			},
			{
				name: "esterror",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Esterror),
			},
			{
				name: "status",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Status),
			},
			{
				name: "constant",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Constant),
			},
			{
				name: "precision",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Precision),
			},
			{
				name: "tolerance",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Tolerance),
			},
			{
				name: "time",
				t:    ArgTypeObject,
				obj:  convertTimeVal(&tx.Time),
			},
			{
				name: "tick",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Tick),
			},
			{
				name: "ppsfreq",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Ppsfreq),
			},
			{
				name: "jitter",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Jitter),
			},
			{
				name: "shift",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Shift),
			},
			{
				name: "stabil",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Stabil),
			},
			{
				name: "jitcnt",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Jitcnt),
			},
			{
				name: "calcnt",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Calcnt),
			},
			{
				name: "errcnt",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Errcnt),
			},
			{
				name: "stbcnt",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Stbcnt),
			},
			{
				name: "tai",
				t:    ArgTypeInt,
				raw:  uintptr(tx.Tai),
			},
		},
	}
}

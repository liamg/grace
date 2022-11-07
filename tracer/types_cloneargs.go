package tracer

import (
	"unsafe"
)

// used by clone3 syscall
type cloneArgs struct {
	Flags      uint64
	PidFd      uint64
	ChildTid   uint64
	ParentTid  uint64
	ExitSignal uint64
	Stack      uint64
	StackSize  uint64
	SetTLS     uint64
	SetTid     uint64
	SetTidSize uint64
	Cgroup     uint64
}

func init() {
	registerTypeHandler(argTypeCloneArgs, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {
			// read the raw C struct from the process memory
			mem, err := readSize(pid, raw, unsafe.Sizeof(cloneArgs{}))
			if err != nil {
				return err
			}

			var args cloneArgs
			if err := decodeStruct(mem, &args); err != nil {
				return err
			}

			arg.obj = convertCloneArgs(args)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertCloneArgs(args cloneArgs) *Object {
	return &Object{
		Name: "clone_args",
		Properties: []Arg{
			{
				name: "flags",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.Flags),
			},
			{
				name: "pid_fd",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.PidFd),
			},
			{
				name: "child_tid",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.ChildTid),
			},
			{
				name: "parent_tid",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.ParentTid),
			},
			{
				name: "exit_signal",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.ExitSignal),
			},
			{
				name: "stack",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.Stack),
			},
			{
				name: "stack_size",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.StackSize),
			},
			{
				name: "set_tls",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.SetTLS),
			},
			{
				name: "set_tid",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.SetTid),
			},
			{
				name: "set_tid_size",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.SetTidSize),
			},
			{
				name: "cgroup",
				t:    ArgTypeUnsignedLong,
				raw:  uintptr(args.Cgroup),
			},
		},
	}
}

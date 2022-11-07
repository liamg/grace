package tracer

import (
	"unsafe"
)

type shmidds struct {
	IPCPerm struct {
		Key     uint32 /* Key provided to shmget */
		Uid     uint32 /* Effective UID of owner */
		Gid     uint32 /* Effective GID of owner */
		Cuid    uint32 /* Effective UID of creator */
		Cgid    uint32 /* Effective GID of creator */
		Mode    uint32 /* Permissions and SHM_DEST + SHM_LOCKED flags */
		Pad1    uint16
		Seq     uint16 /* Sequence */
		Pad2    uint16
		Unused1 uint
		Unused2 uint
	} /* Ownership and permissions */
	Segsz   uint32 /* Size of shared segment (bytes) */
	Atime   uint64 /* Last attach time */
	Dtime   uint64 /* Last detach time */
	Ctime   uint64 /* Last change time */
	Cpid    uint32 /* PID of shared segment creator */
	Lpid    uint32 /* PID of last shmat(2)/shmdt(2) syscall */
	Nattch  uint16 /* Number of current attaches */
	Unused  uint16
	Unused2 uintptr
	Unused3 uintptr
}

func init() {
	registerTypeHandler(argTypeSHMIDDS, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {

		if raw > 0 {

			// read the raw C struct from the process memory
			rawTimeVal, err := readSize(pid, raw, unsafe.Sizeof(shmidds{}))
			if err != nil {
				return err
			}

			var ds shmidds
			if err := decodeStruct(rawTimeVal, &ds); err != nil {
				return err
			}

			arg.obj = convertSHMIDDS(&ds)
			arg.t = ArgTypeObject
		} else {
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}

func convertSHMIDDS(ds *shmidds) *Object {

	var permProps []Arg

	permProps = append(permProps, Arg{
		name: "key",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Key),
	})
	permProps = append(permProps, Arg{
		name: "uid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Uid),
	})
	permProps = append(permProps, Arg{
		name: "gid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Gid),
	})
	permProps = append(permProps, Arg{
		name: "cuid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Cuid),
	})
	permProps = append(permProps, Arg{
		name: "cgid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Cgid),
	})
	permProps = append(permProps, Arg{
		name: "mode",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Mode),
	})
	permProps = append(permProps, Arg{
		name: "seq",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.IPCPerm.Seq),
	})

	var props []Arg

	props = append(props, Arg{
		name: "perm",
		t:    ArgTypeObject,
		obj: &Object{
			Name: "ipc_perm",
		},
	})
	props = append(props, Arg{
		name: "segsize",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.Segsz),
	})

	// TODO: read time structs from process memory instead of storing raw addresses

	props = append(props, Arg{
		name: "atime",
		t:    ArgTypeAddress,
		raw:  uintptr(ds.Atime),
	})
	props = append(props, Arg{
		name: "dtime",
		t:    ArgTypeAddress,
		raw:  uintptr(ds.Dtime),
	})
	props = append(props, Arg{
		name: "ctime",
		t:    ArgTypeAddress,
		raw:  uintptr(ds.Ctime),
	})

	props = append(props, Arg{
		name: "cpid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.Cpid),
	})
	props = append(props, Arg{
		name: "lpid",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.Lpid),
	})
	props = append(props, Arg{
		name: "nattch",
		t:    ArgTypeUnsignedInt,
		raw:  uintptr(ds.Nattch),
	})

	return &Object{
		Name:       "shmid_ds",
		Properties: props,
	}
}

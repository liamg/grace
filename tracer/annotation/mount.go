package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var mountFlags = map[int]string{
	unix.MS_BIND:        "MS_BIND",
	unix.MS_DIRSYNC:     "MS_DIRSYNC",
	unix.MS_MANDLOCK:    "MS_MANDLOCK",
	unix.MS_MOVE:        "MS_MOVE",
	unix.MS_NOATIME:     "MS_NOATIME",
	unix.MS_NODEV:       "MS_NODEV",
	unix.MS_NODIRATIME:  "MS_NODIRATIME",
	unix.MS_NOEXEC:      "MS_NOEXEC",
	unix.MS_NOSUID:      "MS_NOSUID",
	unix.MS_RDONLY:      "MS_RDONLY",
	unix.MS_REC:         "MS_REC",
	unix.MS_RELATIME:    "MS_RELATIME",
	unix.MS_REMOUNT:     "MS_REMOUNT",
	unix.MS_SILENT:      "MS_SILENT",
	unix.MS_STRICTATIME: "MS_STRICTATIME",
	unix.MS_SYNCHRONOUS: "MS_SYNCHRONOUS",
}

func AnnotateMountFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range mountFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var umountFlags = map[int]string{
	unix.MNT_FORCE:       "MNT_FORCE",
	unix.MNT_DETACH:      "MNT_DETACH",
	unix.MNT_EXPIRE:      "MNT_EXPIRE",
	unix.UMOUNT_NOFOLLOW: "UMOUNT_NOFOLLOW",
}

func AnnotateUmountFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range umountFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

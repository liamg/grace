package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var openFlags = map[int]string{
	unix.O_APPEND:    "O_APPEND",
	unix.O_ASYNC:     "O_ASYNC",
	unix.O_CLOEXEC:   "O_CLOEXEC",
	unix.O_CREAT:     "O_CREAT",
	unix.O_DIRECT:    "O_DIRECT",
	unix.O_DIRECTORY: "O_DIRECTORY",
	unix.O_DSYNC:     "O_DSYNC",
	unix.O_EXCL:      "O_EXCL",
	unix.O_NOATIME:   "O_NOATIME",
	unix.O_NOCTTY:    "O_NOCTTY",
	unix.O_NOFOLLOW:  "O_NOFOLLOW",
	unix.O_NONBLOCK:  "O_NONBLOCK",
	unix.O_PATH:      "O_PATH",
	unix.O_SYNC:      "O_SYNC",
	unix.O_TMPFILE:   "O_TMPFILE",
	unix.O_TRUNC:     "O_TRUNC",
	unix.O_RDONLY:    "O_RDONLY",
	unix.O_WRONLY:    "O_WRONLY",
	unix.O_RDWR:      "O_RDWR",
}

func AnnotateOpenFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range openFlags {
		if (int(arg.Raw())&flag != 0) || (int(arg.Raw()) == flag) {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var resolveFlags = map[int]string{
	unix.RESOLVE_BENEATH:       "RESOLVE_BENEATH",
	unix.RESOLVE_IN_ROOT:       "RESOLVE_IN_ROOT",
	unix.RESOLVE_NO_XDEV:       "RESOLVE_NO_XDEV",
	unix.RESOLVE_NO_MAGICLINKS: "RESOLVE_NO_MAGICLINKS",
	unix.RESOLVE_NO_SYMLINKS:   "RESOLVE_NO_SYMLINKS",
	// RESOLVE_CACHED is not available...
}

func AnnotateResolveFlags(arg Arg, _ int) {
	var joins []string
	for flag, name := range resolveFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

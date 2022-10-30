package tracer

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func annotateFd(arg *Arg, pid int) {
	switch int(arg.Raw()) {
	case syscall.Stdin:
		arg.annotation = "stdin"
	case syscall.Stdout:
		arg.annotation = "stdout"
	case syscall.Stderr:
		arg.annotation = "stderr"
	default:
		if path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, arg.Raw())); err == nil {
			arg.annotation = path
		} else if int32(arg.Raw()) == unix.AT_FDCWD {
			arg.annotation = "AT_FDCWD"
			arg.replace = true
		}
	}
}

func annotateAccMode(arg *Arg, _ int) {
	var joins []string
	if arg.Raw() == unix.F_OK {
		joins = append(joins, "F_OK")
	} else {
		if arg.Raw()&unix.R_OK > 0 {
			joins = append(joins, "R_OK")
		}
		if arg.Raw()&unix.W_OK > 0 {
			joins = append(joins, "W_OK")
		}
		if arg.Raw()&unix.X_OK > 0 {
			joins = append(joins, "X_OK")
		}
	}

	arg.annotation = strings.Join(joins, "|")
	arg.replace = arg.annotation != ""
}

func annotateNull(arg *Arg, _ int) {
	if arg.Raw() == 0 {
		arg.annotation = "NULL"
		arg.replace = true
	}
}

func annotateOpenFlags(arg *Arg, pid int) {

	mapFlags := map[int]string{
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
	}

	var joins []string

	for flag, name := range mapFlags {
		if (int(arg.Raw())&flag != 0) || (int(arg.Raw()) == flag) {
			joins = append(joins, name)
		}
	}

	arg.annotation = strings.Join(joins, "|")
	arg.replace = arg.annotation != ""
}

func annotateWhence(arg *Arg, pid int) {
	switch int(arg.Raw()) {
	case unix.SEEK_SET:
		arg.annotation = "SEEK_SET"
	case unix.SEEK_CUR:
		arg.annotation = "SEEK_CUR"
	case unix.SEEK_END:
		arg.annotation = "SEEK_END"
	case unix.SEEK_DATA:
		arg.annotation = "SEEK_DATA"
	case unix.SEEK_HOLE:
		arg.annotation = "SEEK_HOLE"
	}
	arg.replace = arg.annotation != ""
}

func annotateProt(arg *Arg, pid int) {
	if arg.Raw() == unix.PROT_NONE {
		arg.annotation = "PROT_NONE"
		arg.replace = true
		return
	}
	var joins []string
	if arg.Raw()&unix.PROT_READ > 0 {
		joins = append(joins, "PROT_READ")
	}
	if arg.Raw()&unix.PROT_WRITE > 0 {
		joins = append(joins, "PROT_WRITE")
	}
	if arg.Raw()&unix.PROT_EXEC > 0 {
		joins = append(joins, "PROT_EXEC")
	}
	if arg.Raw()&unix.PROT_GROWSUP > 0 {
		joins = append(joins, "PROT_GROWSUP")
	}
	if arg.Raw()&unix.PROT_GROWSDOWN > 0 {
		joins = append(joins, "PROT_GROWSDOWN")
	}

	arg.annotation = strings.Join(joins, "|")
	arg.replace = arg.annotation != ""
}

package annotation

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/liamg/grace/tracer/netw"

	"golang.org/x/sys/unix"
)

func AnnotateFd(arg Arg, pid int) {

	if path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, arg.Raw())); err == nil {

		switch {
		case strings.HasPrefix(path, "socket:["):
			fd := strings.TrimPrefix(path, "socket:[")
			fd = strings.TrimSuffix(fd, "]")
			if inode, err := strconv.Atoi(fd); err == nil {
				if str := socketInoToString(inode); str != "" {
					arg.SetAnnotation(str, false)
					return
				}
			}
		}
		arg.SetAnnotation(path, false)
	}

	if int32(arg.Raw()) == unix.AT_FDCWD {
		arg.SetAnnotation("AT_FDCWD", true)
		return
	}

	switch int(arg.Raw()) {
	case syscall.Stdin:
		arg.SetAnnotation("STDIN", false)
	case syscall.Stdout:
		arg.SetAnnotation("STDOUT", false)
	case syscall.Stderr:
		arg.SetAnnotation("STDERR", false)
	}
}

func socketInoToString(inode int) string {
	if conns, err := netw.ListConnections(); err == nil {
		for _, conn := range conns {
			if conn.INode == inode {
				return fmt.Sprintf("%s %s:%d -> %s:%d", conn.Protocol, conn.LocalAddress, conn.LocalPort, conn.RemoteAddress, conn.RemotePort)
			}
		}
	}
	return ""
}

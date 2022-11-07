package annotation

import "golang.org/x/sys/unix"

var epollCtlOps = map[int]string{
	unix.EPOLL_CTL_ADD: "EPOLL_CTL_ADD",
	unix.EPOLL_CTL_DEL: "EPOLL_CTL_DEL",
	unix.EPOLL_CTL_MOD: "EPOLL_CTL_MOD",
}

func AnnotateEpollCtlOp(arg Arg, _ int) {
	if str, ok := epollCtlOps[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

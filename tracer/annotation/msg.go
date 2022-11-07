package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var msgFlags = map[int]string{
	unix.MSG_OOB:          "MSG_OOB",
	unix.MSG_PEEK:         "MSG_PEEK",
	unix.MSG_DONTROUTE:    "MSG_DONTROUTE",
	unix.MSG_CTRUNC:       "MSG_CTRUNC",
	unix.MSG_PROXY:        "MSG_PROXY",
	unix.MSG_TRUNC:        "MSG_TRUNC",
	unix.MSG_DONTWAIT:     "MSG_DONTWAIT",
	unix.MSG_EOR:          "MSG_EOR",
	unix.MSG_WAITALL:      "MSG_WAITALL",
	unix.MSG_FIN:          "MSG_FIN",
	unix.MSG_SYN:          "MSG_SYN",
	unix.MSG_CONFIRM:      "MSG_CONFIRM",
	unix.MSG_RST:          "MSG_RST",
	unix.MSG_ERRQUEUE:     "MSG_ERRQUEUE",
	unix.MSG_NOSIGNAL:     "MSG_NOSIGNAL",
	unix.MSG_MORE:         "MSG_MORE",
	unix.MSG_WAITFORONE:   "MSG_WAITFORONE",
	unix.MSG_FASTOPEN:     "MSG_FASTOPEN",
	unix.MSG_CMSG_CLOEXEC: "MSG_CMSG_CLOEXEC",
}

func AnnotateMsgFlags(arg Arg, _ int) {
	var joins []string
	for flag, str := range msgFlags {
		if arg.Raw()&uintptr(flag) != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

func AnnotateMsgType(arg Arg, pid int) {
	// TODO:
}

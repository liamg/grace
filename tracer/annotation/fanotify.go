package annotation

import (
	"strings"

	"golang.org/x/sys/unix"
)

var fanotifyFlags = map[int]string{
	unix.FAN_CLASS_PRE_CONTENT: "FAN_CLASS_PRE_CONTENT",
	unix.FAN_CLASS_CONTENT:     "FAN_CLASS_CONTENT",
	unix.FAN_CLASS_NOTIF:       "FAN_CLASS_NOTIF",
	unix.FAN_REPORT_FID:        "FAN_REPORT_FID",
	unix.FAN_REPORT_NAME:       "FAN_REPORT_NAME",
	unix.FAN_REPORT_DIR_FID:    "FAN_REPORT_DIR_FID",
	unix.FAN_REPORT_TID:        "FAN_REPORT_TID",
	unix.FAN_NONBLOCK:          "FAN_NONBLOCK",
	unix.FAN_UNLIMITED_QUEUE:   "FAN_UNLIMITED_QUEUE",
	unix.FAN_UNLIMITED_MARKS:   "FAN_UNLIMITED_MARKS",
	unix.FAN_ENABLE_AUDIT:      "FAN_ENABLE_AUDIT",
}

var fanotifyMarkFlags = map[int]string{
	unix.FAN_MARK_ADD:                 "FAN_MARK_ADD",
	unix.FAN_MARK_REMOVE:              "FAN_MARK_REMOVE",
	unix.FAN_MARK_DONT_FOLLOW:         "FAN_MARK_DONT_FOLLOW",
	unix.FAN_MARK_ONLYDIR:             "FAN_MARK_ONLYDIR",
	unix.FAN_MARK_MOUNT:               "FAN_MARK_MOUNT",
	unix.FAN_MARK_IGNORED_MASK:        "FAN_MARK_IGNORED_MASK",
	unix.FAN_MARK_IGNORED_SURV_MODIFY: "FAN_MARK_IGNORED_SURV_MODIFY",
	unix.FAN_MARK_FLUSH:               "FAN_MARK_FLUSH",
	unix.FAN_MARK_FILESYSTEM:          "FAN_MARK_FILESYSTEM",
}

func AnnotateFANotifyFlags(arg Arg, pid int) {
	var joins []string
	for flag, str := range fanotifyFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

func AnnotateFANotifyEventFlags(arg Arg, pid int) {
	AnnotateOpenFlags(arg, pid)
}

func AnnotateFANotifyMarkFlags(arg Arg, pid int) {
	var joins []string
	for flag, str := range fanotifyMarkFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, str)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

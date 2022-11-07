package annotation

import "golang.org/x/sys/unix"

var syslogTypes = map[int]string{
	unix.SYSLOG_ACTION_READ:          "SYSLOG_ACTION_READ",
	unix.SYSLOG_ACTION_READ_ALL:      "SYSLOG_ACTION_READ_ALL",
	unix.SYSLOG_ACTION_READ_CLEAR:    "SYSLOG_ACTION_READ_CLEAR",
	unix.SYSLOG_ACTION_CLEAR:         "SYSLOG_ACTION_CLEAR",
	unix.SYSLOG_ACTION_CONSOLE_OFF:   "SYSLOG_ACTION_CONSOLE_OFF",
	unix.SYSLOG_ACTION_CONSOLE_ON:    "SYSLOG_ACTION_CONSOLE_ON",
	unix.SYSLOG_ACTION_CONSOLE_LEVEL: "SYSLOG_ACTION_CONSOLE_LEVEL",
	unix.SYSLOG_ACTION_SIZE_UNREAD:   "SYSLOG_ACTION_SIZE_UNREAD",
	unix.SYSLOG_ACTION_SIZE_BUFFER:   "SYSLOG_ACTION_SIZE_BUFFER",
}

func AnnotateSyslogType(arg Arg, _ int) {
	if name, ok := syslogTypes[int(arg.Raw())]; ok {
		arg.SetAnnotation(name, true)
	}
}

package annotation

import "golang.org/x/sys/unix"

var keyctlCmds = map[int]string{
	unix.KEYCTL_GET_KEYRING_ID:       "KEYCTL_GET_KEYRING_ID",
	unix.KEYCTL_JOIN_SESSION_KEYRING: "KEYCTL_JOIN_SESSION_KEYRING",
	unix.KEYCTL_UPDATE:               "KEYCTL_UPDATE",
	unix.KEYCTL_REVOKE:               "KEYCTL_REVOKE",
	unix.KEYCTL_DESCRIBE:             "KEYCTL_DESCRIBE",
	unix.KEYCTL_CLEAR:                "KEYCTL_CLEAR",
	unix.KEYCTL_LINK:                 "KEYCTL_LINK",
	unix.KEYCTL_UNLINK:               "KEYCTL_UNLINK",
	unix.KEYCTL_SEARCH:               "KEYCTL_SEARCH",
	unix.KEYCTL_READ:                 "KEYCTL_READ",
	unix.KEYCTL_CHOWN:                "KEYCTL_CHOWN",
	unix.KEYCTL_SETPERM:              "KEYCTL_SETPERM",
	unix.KEYCTL_INSTANTIATE:          "KEYCTL_INSTANTIATE",
	unix.KEYCTL_NEGATE:               "KEYCTL_NEGATE",
	unix.KEYCTL_SET_REQKEY_KEYRING:   "KEYCTL_SET_REQKEY_KEYRING",
	unix.KEYCTL_SET_TIMEOUT:          "KEYCTL_SET_TIMEOUT",
	unix.KEYCTL_ASSUME_AUTHORITY:     "KEYCTL_ASSUME_AUTHORITY",
	unix.KEYCTL_GET_SECURITY:         "KEYCTL_GET_SECURITY",
	unix.KEYCTL_SESSION_TO_PARENT:    "KEYCTL_SESSION_TO_PARENT",
	unix.KEYCTL_REJECT:               "KEYCTL_REJECT",
	unix.KEYCTL_INSTANTIATE_IOV:      "KEYCTL_INSTANTIATE_IOV",
	unix.KEYCTL_INVALIDATE:           "KEYCTL_INVALIDATE",
}

func AnnotateKeyctlCommand(arg Arg, pid int) {
	if cmd, ok := keyctlCmds[int(arg.Raw())]; ok {
		arg.SetAnnotation(cmd, true)
	}
}

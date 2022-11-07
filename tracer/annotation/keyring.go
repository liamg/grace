package annotation

import "golang.org/x/sys/unix"

var keyringIDs = map[int]string{
	unix.KEY_SPEC_THREAD_KEYRING:       "KEY_SPEC_THREAD_KEYRING",
	unix.KEY_SPEC_PROCESS_KEYRING:      "KEY_SPEC_PROCESS_KEYRING",
	unix.KEY_SPEC_SESSION_KEYRING:      "KEY_SPEC_SESSION_KEYRING",
	unix.KEY_SPEC_USER_KEYRING:         "KEY_SPEC_USER_KEYRING",
	unix.KEY_SPEC_USER_SESSION_KEYRING: "KEY_SPEC_USER_SESSION_KEYRING",
	unix.KEY_SPEC_GROUP_KEYRING:        "KEY_SPEC_GROUP_KEYRING",
}

func AnnotateKeyringID(arg Arg, _ int) {
	if id, ok := keyringIDs[int(arg.Raw())]; ok {
		arg.SetAnnotation(id, true)
	}
}

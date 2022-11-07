package annotation

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func AnnotateDevice(arg Arg, pid int) {
	arg.SetAnnotation(DeviceToString(uint64(arg.Raw())), true)
}

func DeviceToString(dev uint64) string {
	major := "0"
	if m := unix.Major(dev); m != 0 {
		major = fmt.Sprintf("0x%x", m)
	}
	minor := "0"
	if m := unix.Minor(dev); m != 0 {
		minor = fmt.Sprintf("0x%x", m)
	}
	return fmt.Sprintf("makedev(%s, %s)", major, minor)
}

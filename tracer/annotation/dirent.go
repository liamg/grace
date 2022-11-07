package annotation

import "golang.org/x/sys/unix"

var dirEntTypes = map[int]string{
	unix.DT_BLK:     "DT_BLK",
	unix.DT_CHR:     "DT_CHR",
	unix.DT_DIR:     "DT_DIR",
	unix.DT_FIFO:    "DT_FIFO",
	unix.DT_LNK:     "DT_LNK",
	unix.DT_REG:     "DT_REG",
	unix.DT_SOCK:    "DT_SOCK",
	unix.DT_WHT:     "DT_WHT",
	unix.DT_UNKNOWN: "DT_UNKNOWN",
}

func AnnotateDirEntType(arg Arg, _ int) {
	if str, ok := dirEntTypes[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

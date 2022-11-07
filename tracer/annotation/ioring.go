package annotation

import "strings"

var ioringEnterFlags = map[int]string{
	1:  "IORING_ENTER_GETEVENTS",
	2:  "IORING_ENTER_SQ_WAKEUP",
	4:  "IORING_ENTER_SQ_WAIT",
	8:  "IORING_ENTER_EXT_ARG",
	16: "IORING_ENTER_REGISTERED_RING",
}

func AnnotateIoUringEnterFlags(arg Arg, pid int) {
	var joins []string
	for flag, name := range ioringEnterFlags {
		if int(arg.Raw())&flag != 0 {
			joins = append(joins, name)
		}
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var ioringOpCodes = []string{
	"IORING_REGISTER_BUFFERS",
	"IORING_UNREGISTER_BUFFERS",
	"IORING_REGISTER_FILES",
	"IORING_UNREGISTER_FILES",
	"IORING_REGISTER_EVENTFD",
	"IORING_UNREGISTER_EVENTFD",
	"IORING_REGISTER_FILES_UPDATE",
	"IORING_REGISTER_EVENTFD_ASYNC",
	"IORING_REGISTER_PROBE",
	"IORING_REGISTER_PERSONALITY",
	"IORING_UNREGISTER_PERSONALITY",
	"IORING_REGISTER_RESTRICTIONS",
	"IORING_REGISTER_ENABLE_RINGS",
	"IORING_REGISTER_FILES2",
	"IORING_REGISTER_FILES_UPDATE2",
	"IORING_REGISTER_BUFFERS2",
	"IORING_REGISTER_BUFFERS_UPDATE",
	"IORING_REGISTER_IOWQ_AFF",
	"IORING_UNREGISTER_IOWQ_AFF",
	"IORING_REGISTER_IOWQ_MAX_WORKERS",
	"IORING_REGISTER_RING_FDS",
	"IORING_UNREGISTER_RING_FDS",
	"IORING_REGISTER_PBUF_RING",
	"IORING_UNREGISTER_PBUF_RING",
	"IORING_REGISTER_SYNC_CANCEL",
	"IORING_REGISTER_FILE_ALLOC_RANGE",
}

func AnnotateIORingOpCode(arg Arg, _ int) {
	opcode := int(arg.Raw())
	if opcode >= 0 && opcode < len(ioringOpCodes) {
		arg.SetAnnotation(ioringOpCodes[opcode], true)
	}
}

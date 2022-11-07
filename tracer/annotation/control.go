package annotation

import "golang.org/x/sys/unix"

var socketLevels = map[int]string{
	unix.SOL_SOCKET:   "SOL_SOCKET",
	unix.SOL_AAL:      "SOL_AAL",
	unix.SOL_ALG:      "SOL_ALG",
	unix.SOL_ATM:      "SOL_ATM",
	unix.SOL_CAIF:     "SOL_CAIF",
	unix.SOL_CAN_BASE: "SOL_CAN_BASE",
	unix.SOL_CAN_RAW:  "SOL_CAN_RAW",
	unix.SOL_DCCP:     "SOL_DCCP",
	unix.SOL_DECNET:   "SOL_DECNET",
	unix.SOL_ICMPV6:   "SOL_ICMPV6",
	unix.SOL_IP:       "SOL_IP",
	unix.SOL_IPV6:     "SOL_IPV6",
	unix.SOL_IRDA:     "SOL_IRDA",
	unix.SOL_IUCV:     "SOL_IUCV",
	unix.SOL_KCM:      "SOL_KCM",
	unix.SOL_LLC:      "SOL_LLC",
	unix.SOL_MCTP:     "SOL_MCTP",
	unix.SOL_MPTCP:    "SOL_MPTCP",
	unix.SOL_NETBEUI:  "SOL_NETBEUI",
	unix.SOL_NETLINK:  "SOL_NETLINK",
	unix.SOL_NFC:      "SOL_NFC",
	unix.SOL_PACKET:   "SOL_PACKET",
	unix.SOL_PNPIPE:   "SOL_PNPIPE",
	unix.SOL_PPPOL2TP: "SOL_PPPOL2TP",
	unix.SOL_RAW:      "SOL_RAW",
	unix.SOL_RDS:      "SOL_RDS",
	unix.SOL_RXRPC:    "SOL_RXRPC",
	unix.SOL_SMC:      "SOL_SMC",
	unix.SOL_TCP:      "SOL_TCP",
	unix.SOL_TIPC:     "SOL_TIPC",
	unix.SOL_TLS:      "SOL_TLS",
	unix.SOL_X25:      "SOL_X25",
	unix.SOL_XDP:      "SOL_XDP",
}

func AnnotateSocketLevel(arg Arg, _ int) {
	if str, ok := socketLevels[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var cmsgTypes = map[int]string{
	unix.SCM_RIGHTS:      "SCM_RIGHTS",
	unix.SCM_CREDENTIALS: "SCM_CREDENTIALS",
	unix.SCM_TIMESTAMP:   "SCM_TIMESTAMP",
}

func AnnotateControlMessageType(arg Arg, _ int) {
	if str, ok := cmsgTypes[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

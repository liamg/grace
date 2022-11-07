package annotation

import "golang.org/x/sys/unix"

var socketDomains = map[int]string{
	unix.AF_UNIX:      "AF_UNIX",
	unix.AF_INET:      "AF_INET",
	unix.AF_INET6:     "AF_INET6",
	unix.AF_AX25:      "AF_AX25",
	unix.AF_IPX:       "AF_IPX",
	unix.AF_APPLETALK: "AF_APPLETALK",
	unix.AF_X25:       "AF_X25",
	unix.AF_DECnet:    "AF_DECnet",
	unix.AF_KEY:       "AF_KEY",
	unix.AF_NETLINK:   "AF_NETLINK",
	unix.AF_PACKET:    "AF_PACKET",
	unix.AF_RDS:       "AF_RDS",
	unix.AF_PPPOX:     "AF_PPPOX",
	unix.AF_LLC:       "AF_LLC",
	unix.AF_IB:        "AF_IB",
	unix.AF_MPLS:      "AF_MPLS",
	unix.AF_CAN:       "AF_CAN",
	unix.AF_TIPC:      "AF_TIPC",
	unix.AF_BLUETOOTH: "AF_BLUETOOTH",
	unix.AF_ALG:       "AF_ALG",
	unix.AF_VSOCK:     "AF_VSOCK",
	unix.AF_KCM:       "AF_KCM",
	unix.AF_XDP:       "AF_XDP",
}

func SocketFamily(family int) string {
	return socketDomains[family]
}

func AnnotateSocketDomain(arg Arg, _ int) {
	if str, ok := socketDomains[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var socketTypes = map[int]string{
	unix.SOCK_STREAM:    "SOCK_STREAM",
	unix.SOCK_DGRAM:     "SOCK_DGRAM",
	unix.SOCK_SEQPACKET: "SOCK_SEQPACKET",
	unix.SOCK_RAW:       "SOCK_RAW",
	unix.SOCK_RDM:       "SOCK_RDM",
	unix.SOCK_PACKET:    "SOCK_PACKET",
	unix.SOCK_NONBLOCK:  "SOCK_NONBLOCK",
	unix.SOCK_CLOEXEC:   "SOCK_CLOEXEC",
}

func AnnotateSocketType(arg Arg, _ int) {
	if str, ok := socketTypes[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var socketProtocols = map[int]string{
	unix.IPPROTO_IP:      "IPPROTO_IP",
	unix.IPPROTO_IPV6:    "IPPROTO_IPV6",
	unix.IPPROTO_ICMP:    "IPPROTO_ICMP",
	unix.IPPROTO_ICMPV6:  "IPPROTO_ICMPV6",
	unix.IPPROTO_IGMP:    "IPPROTO_IGMP",
	unix.IPPROTO_IPIP:    "IPPROTO_IPIP",
	unix.IPPROTO_TCP:     "IPPROTO_TCP",
	unix.IPPROTO_EGP:     "IPPROTO_EGP",
	unix.IPPROTO_PUP:     "IPPROTO_PUP",
	unix.IPPROTO_UDP:     "IPPROTO_UDP",
	unix.IPPROTO_IDP:     "IPPROTO_IDP",
	unix.IPPROTO_TP:      "IPPROTO_TP",
	unix.IPPROTO_DCCP:    "IPPROTO_DCCP",
	unix.IPPROTO_RSVP:    "IPPROTO_RSVP",
	unix.IPPROTO_GRE:     "IPPROTO_GRE",
	unix.IPPROTO_ESP:     "IPPROTO_ESP",
	unix.IPPROTO_AH:      "IPPROTO_AH",
	unix.IPPROTO_MTP:     "IPPROTO_MTP",
	unix.IPPROTO_BEETPH:  "IPPROTO_BEETPH",
	unix.IPPROTO_ENCAP:   "IPPROTO_ENCAP",
	unix.IPPROTO_PIM:     "IPPROTO_PIM",
	unix.IPPROTO_COMP:    "IPPROTO_COMP",
	unix.IPPROTO_SCTP:    "IPPROTO_SCTP",
	unix.IPPROTO_UDPLITE: "IPPROTO_UDPLITE",
	unix.IPPROTO_MPLS:    "IPPROTO_MPLS",
	unix.IPPROTO_RAW:     "IPPROTO_RAW",
}

func AnnotateSocketProtocol(arg Arg, _ int) {
	if str, ok := socketProtocols[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

var socketOptions = map[int]string{
	unix.SO_ACCEPTCONN:                    "SO_ACCEPTCONN",
	unix.SO_ATTACH_BPF:                    "SO_ATTACH_BPF",
	unix.SO_ATTACH_REUSEPORT_CBPF:         "SO_ATTACH_REUSEPORT_CBPF",
	unix.SO_ATTACH_REUSEPORT_EBPF:         "SO_ATTACH_REUSEPORT_EBPF",
	unix.SO_BINDTODEVICE:                  "SO_BINDTODEVICE",
	unix.SO_BINDTOIFINDEX:                 "SO_BINDTOIFINDEX",
	unix.SO_BPF_EXTENSIONS:                "SO_BPF_EXTENSIONS",
	unix.SO_BROADCAST:                     "SO_BROADCAST",
	unix.SO_BSDCOMPAT:                     "SO_BSDCOMPAT",
	unix.SO_BUF_LOCK:                      "SO_BUF_LOCK",
	unix.SO_BUSY_POLL:                     "SO_BUSY_POLL",
	unix.SO_BUSY_POLL_BUDGET:              "SO_BUSY_POLL_BUDGET",
	unix.SO_CNX_ADVICE:                    "SO_CNX_ADVICE",
	unix.SO_COOKIE:                        "SO_COOKIE",
	unix.SO_DETACH_REUSEPORT_BPF:          "SO_DETACH_REUSEPORT_BPF",
	unix.SO_DOMAIN:                        "SO_DOMAIN",
	unix.SO_DONTROUTE:                     "SO_DONTROUTE",
	unix.SO_ERROR:                         "SO_ERROR",
	unix.SO_INCOMING_CPU:                  "SO_INCOMING_CPU",
	unix.SO_INCOMING_NAPI_ID:              "SO_INCOMING_NAPI_ID",
	unix.SO_KEEPALIVE:                     "SO_KEEPALIVE",
	unix.SO_LINGER:                        "SO_LINGER",
	unix.SO_LOCK_FILTER:                   "SO_LOCK_FILTER",
	unix.SO_MARK:                          "SO_MARK",
	unix.SO_MAX_PACING_RATE:               "SO_MAX_PACING_RATE",
	unix.SO_MEMINFO:                       "SO_MEMINFO",
	unix.SO_NETNS_COOKIE:                  "SO_NETNS_COOKIE",
	unix.SO_NOFCS:                         "SO_NOFCS",
	unix.SO_OOBINLINE:                     "SO_OOBINLINE",
	unix.SO_PASSCRED:                      "SO_PASSCRED",
	unix.SO_PASSSEC:                       "SO_PASSSEC",
	unix.SO_PEEK_OFF:                      "SO_PEEK_OFF",
	unix.SO_PEERCRED:                      "SO_PEERCRED",
	unix.SO_PEERGROUPS:                    "SO_PEERGROUPS",
	unix.SO_PEERSEC:                       "SO_PEERSEC",
	unix.SO_PREFER_BUSY_POLL:              "SO_PREFER_BUSY_POLL",
	unix.SO_PROTOCOL:                      "SO_PROTOCOL",
	unix.SO_RCVBUF:                        "SO_RCVBUF",
	unix.SO_RCVBUFFORCE:                   "SO_RCVBUFFORCE",
	unix.SO_RCVLOWAT:                      "SO_RCVLOWAT",
	unix.SO_RCVMARK:                       "SO_RCVMARK",
	unix.SO_RCVTIMEO_NEW:                  "SO_RCVTIMEO_NEW",
	unix.SO_RCVTIMEO_OLD:                  "SO_RCVTIMEO_OLD",
	unix.SO_RESERVE_MEM:                   "SO_RESERVE_MEM",
	unix.SO_REUSEADDR:                     "SO_REUSEADDR",
	unix.SO_REUSEPORT:                     "SO_REUSEPORT",
	unix.SO_RXQ_OVFL:                      "SO_RXQ_OVFL",
	unix.SO_SECURITY_AUTHENTICATION:       "SO_SECURITY_AUTHENTICATION",
	unix.SO_SECURITY_ENCRYPTION_NETWORK:   "SO_SECURITY_ENCRYPTION_NETWORK",
	unix.SO_SECURITY_ENCRYPTION_TRANSPORT: "SO_SECURITY_ENCRYPTION_TRANSPORT",
	unix.SO_SELECT_ERR_QUEUE:              "SO_SELECT_ERR_QUEUE",
	unix.SO_SNDBUF:                        "SO_SNDBUF",
	unix.SO_SNDBUFFORCE:                   "SO_SNDBUFFORCE",
	unix.SO_SNDLOWAT:                      "SO_SNDLOWAT",
	unix.SO_SNDTIMEO_NEW:                  "SO_SNDTIMEO_NEW",
	unix.SO_SNDTIMEO_OLD:                  "SO_SNDTIMEO_OLD",
	unix.SO_TIMESTAMPING_NEW:              "SO_TIMESTAMPING_NEW",
	unix.SO_TIMESTAMPING_OLD:              "SO_TIMESTAMPING_OLD",
	unix.SO_TIMESTAMPNS_NEW:               "SO_TIMESTAMPNS_NEW",
	unix.SO_TIMESTAMPNS_OLD:               "SO_TIMESTAMPNS_OLD",
	unix.SO_TIMESTAMP_NEW:                 "SO_TIMESTAMP_NEW",
	unix.SO_TXREHASH:                      "SO_TXREHASH",
	unix.SO_TXTIME:                        "SO_TXTIME",
	unix.SO_TYPE:                          "SO_TYPE",
	unix.SO_WIFI_STATUS:                   "SO_WIFI_STATUS",
	unix.SO_ZEROCOPY:                      "SO_ZEROCOPY",
	unix.SO_PEERNAME:                      "SO_PEERNAME",
	unix.SO_PRIORITY:                      "SO_PRIORITY",
	unix.SO_TIMESTAMP_OLD:                 "SO_TIMESTAMP_OLD",
	unix.SO_DETACH_FILTER:                 "SO_DETACH_FILTER",
	unix.SO_EE_ORIGIN_NONE:                "SO_EE_ORIGIN_NONE",
	unix.SO_EE_RFC4884_FLAG_INVALID:       "SO_EE_RFC4884_FLAG_INVALID",
	unix.SO_GET_FILTER:                    "SO_GET_FILTER",
	unix.SO_NO_CHECK:                      "SO_NO_CHECK",
}

func AnnotateSocketOption(arg Arg, _ int) {
	if str, ok := socketOptions[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

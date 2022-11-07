package annotation

import "C"

/*
   #include<errno.h>

   // the following definitions are never meant to be seen by user programs, other than via ptrace...
   #ifndef ERESTARTSYS
   #define ERESTARTSYS 512
   #endif

   #ifndef ERESTARTNOINTR
   #define ERESTARTNOINTR	 513
   #endif

   #ifndef ERESTARTNOHAND
   #define ERESTARTNOHAND 514
   #endif

   #ifndef  ENOIOCTLCMD
   #define  ENOIOCTLCMD 515
   #endif

   #ifndef ERESTART_RESTARTBLOCK
   #define ERESTART_RESTARTBLOCK 516
   #endif

   // TODO: add NFSv3 errors

*/
import "C"

import (
	"syscall"
)

var errors = map[int]string{
	C.ERESTARTSYS:           "ERESTARTSYS",
	C.ERESTARTNOINTR:        "ERESTARTNOINTR",
	C.ERESTARTNOHAND:        "ERESTARTNOHAND",
	C.ENOIOCTLCMD:           "ENOIOCTLCMD",
	C.ERESTART_RESTARTBLOCK: "ERESTART_RESTARTBLOCK",
	C.ERESTART:              "ERESTART",
	C.EPERM:                 "EPERM",
	C.ENOENT:                "ENOENT",
	C.ESRCH:                 "ESRCH",
	C.EINTR:                 "EINTR",
	C.EIO:                   "EIO",
	C.ENXIO:                 "ENXIO",
	C.E2BIG:                 "E2BIG",
	C.ENOEXEC:               "ENOEXEC",
	C.EBADF:                 "EBADF",
	C.ECHILD:                "ECHILD",
	C.EAGAIN:                "EAGAIN",
	C.ENOMEM:                "ENOMEM",
	C.EACCES:                "EACCES",
	C.EFAULT:                "EFAULT",
	C.ENOTBLK:               "ENOTBLK",
	C.EBUSY:                 "EBUSY",
	C.EEXIST:                "EEXIST",
	C.EXDEV:                 "EXDEV",
	C.ENODEV:                "ENODEV",
	C.ENOTDIR:               "ENOTDIR",
	C.EISDIR:                "EISDIR",
	C.EINVAL:                "EINVAL",
	C.ENFILE:                "ENFILE",
	C.EMFILE:                "EMFILE",
	C.ENOTTY:                "ENOTTY",
	C.ETXTBSY:               "ETXTBSY",
	C.EFBIG:                 "EFBIG",
	C.ENOSPC:                "ENOSPC",
	C.ESPIPE:                "ESPIPE",
	C.EROFS:                 "EROFS",
	C.EMLINK:                "EMLINK",
	C.EPIPE:                 "EPIPE",
	C.EDOM:                  "EDOM",
	C.ERANGE:                "ERANGE",
	C.EDEADLK:               "EDEADLK",
	C.ENAMETOOLONG:          "ENAMETOOLONG",
	C.ENOLCK:                "ENOLCK",
	C.ENOSYS:                "ENOSYS",
	C.ENOTEMPTY:             "ENOTEMPTY",
	C.ELOOP:                 "ELOOP",
	C.ENOMSG:                "ENOMSG",
	C.EIDRM:                 "EIDRM",
	C.ECHRNG:                "ECHRNG",
	C.EL2NSYNC:              "EL2NSYNC",
	C.EL3HLT:                "EL3HLT",
	C.EL3RST:                "EL3RST",
	C.ELNRNG:                "ELNRNG",
	C.EUNATCH:               "EUNATCH",
	C.ENOCSI:                "ENOCSI",
	C.EL2HLT:                "EL2HLT",
	C.EBADE:                 "EBADE",
	C.EBADR:                 "EBADR",
	C.EXFULL:                "EXFULL",
	C.ENOANO:                "ENOANO",
	C.EBADRQC:               "EBADRQC",
	C.EBADSLT:               "EBADSLT",
	C.EBFONT:                "EBFONT",
	C.ENOSTR:                "ENOSTR",
	C.ENODATA:               "ENODATA",
	C.ETIME:                 "ETIME",
	C.ENOSR:                 "ENOSR",
	C.ENONET:                "ENONET",
	C.ENOPKG:                "ENOPKG",
	C.EREMOTE:               "EREMOTE",
	C.ENOLINK:               "ENOLINK",
	C.EADV:                  "EADV",
	C.ESRMNT:                "ESRMNT",
	C.ECOMM:                 "ECOMM",
	C.EPROTO:                "EPROTO",
	C.EMULTIHOP:             "EMULTIHOP",
	C.EBADMSG:               "EBADMSG",
	C.ENOTUNIQ:              "ENOTUNIQ",
	C.EBADFD:                "EBADFD",
	C.EREMCHG:               "EREMCHG",
	C.ELIBACC:               "ELIBACC",
	C.ELIBBAD:               "ELIBBAD",
	C.ELIBSCN:               "ELIBSCN",
	C.ELIBMAX:               "ELIBMAX",
	C.ELIBEXEC:              "ELIBEXEC",
	C.EILSEQ:                "EILSEQ",
	C.ESTRPIPE:              "ESTRPIPE",
	C.EUSERS:                "EUSERS",
	C.ENOTSOCK:              "ENOTSOCK",
	C.EDESTADDRREQ:          "EDESTADDRREQ",
	C.EMSGSIZE:              "EMSGSIZE",
	C.EPROTOTYPE:            "EPROTOTYPE",
	C.ENOPROTOOPT:           "ENOPROTOOPT",
	C.EPROTONOSUPPORT:       "EPROTONOSUPPORT",
	C.ESOCKTNOSUPPORT:       "ESOCKTNOSUPPORT",
	C.EOPNOTSUPP:            "EOPNOTSUPP",
	C.EPFNOSUPPORT:          "EPFNOSUPPORT",
	C.EAFNOSUPPORT:          "EAFNOSUPPORT",
	C.EADDRINUSE:            "EADDRINUSE",
	C.EADDRNOTAVAIL:         "EADDRNOTAVAIL",
	C.ENETDOWN:              "ENETDOWN",
	C.ENETUNREACH:           "ENETUNREACH",
	C.ENETRESET:             "ENETRESET",
	C.ECONNABORTED:          "ECONNABORTED",
	C.ECONNRESET:            "ECONNRESET",
	C.ENOBUFS:               "ENOBUFS",
	C.EISCONN:               "EISCONN",
	C.ENOTCONN:              "ENOTCONN",
	C.ESHUTDOWN:             "ESHUTDOWN",
	C.ETOOMANYREFS:          "ETOOMANYREFS",
	C.ETIMEDOUT:             "ETIMEDOUT",
	C.ECONNREFUSED:          "ECONNREFUSED",
	C.EHOSTDOWN:             "EHOSTDOWN",
	C.EHOSTUNREACH:          "EHOSTUNREACH",
	C.EALREADY:              "EALREADY",
	C.EINPROGRESS:           "EINPROGRESS",
	C.ESTALE:                "ESTALE",
	C.EUCLEAN:               "EUCLEAN",
	C.ENOTNAM:               "ENOTNAM",
	C.ENAVAIL:               "ENAVAIL",
	C.EISNAM:                "EISNAM",
	C.EREMOTEIO:             "EREMOTEIO",
	C.EDQUOT:                "EDQUOT",
	C.ENOMEDIUM:             "ENOMEDIUM",
	C.EMEDIUMTYPE:           "EMEDIUMTYPE",
	C.ECANCELED:             "ECANCELED",
	C.ENOKEY:                "ENOKEY",
	C.EKEYEXPIRED:           "EKEYEXPIRED",
	C.EKEYREVOKED:           "EKEYREVOKED",
	C.EKEYREJECTED:          "EKEYREJECTED",
	C.EOWNERDEAD:            "EOWNERDEAD",
	C.ENOTRECOVERABLE:       "ENOTRECOVERABLE",
	C.ERFKILL:               "ERFKILL",
	C.EHWPOISON:             "EHWPOISON",
}

func ErrNoToString(errno int) string {
	if err, ok := errors[errno]; ok {
		return err
	}
	return syscall.Errno(errno).Error()
}

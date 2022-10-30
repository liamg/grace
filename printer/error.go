package printer

import (
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/liamg/grace/tracer"
)

func (p *Printer) printError(colour Colour, arg *tracer.Arg) {

	code := -arg.Int()
	if code == 0 {
		p.PrintColour(ColourGreen, "0")
		return
	}

	constant := "unknown"

	switch syscall.Errno(code) {
	case unix.E2BIG:
		constant = "E2BIG"
	case unix.EACCES:
		constant = "EACCES"
	case unix.EADDRINUSE:
		constant = "EADDRINUSE"
	case unix.EADDRNOTAVAIL:
		constant = "EADDRNOTAVAIL"
	case unix.EADV:
		constant = "EADV"
	case unix.EAFNOSUPPORT:
		constant = "EAFNOSUPPORT"
	case unix.EAGAIN:
		constant = "EAGAIN"
	case unix.EALREADY:
		constant = "EALREADY"
	case unix.EBADE:
		constant = "EBADE"
	case unix.EBADF:
		constant = "EBADF"
	case unix.EBADFD:
		constant = "EBADFD"
	case unix.EBADMSG:
		constant = "EBADMSG"
	case unix.EBADR:
		constant = "EBADR"
	case unix.EBADRQC:
		constant = "EBADRQC"
	case unix.EBADSLT:
		constant = "EBADSLT"
	case unix.EBFONT:
		constant = "EBFONT"
	case unix.EBUSY:
		constant = "EBUSY"
	case unix.ECANCELED:
		constant = "ECANCELED"
	case unix.ECHILD:
		constant = "ECHILD"
	case unix.ECHRNG:
		constant = "ECHRNG"
	case unix.ECOMM:
		constant = "ECOMM"
	case unix.ECONNABORTED:
		constant = "ECONNABORTED"
	case unix.ECONNREFUSED:
		constant = "ECONNREFUSED"
	case unix.ECONNRESET:
		constant = "ECONNRESET"
	case unix.EDEADLOCK:
		constant = "EDEADLOCK"
	case unix.EDESTADDRREQ:
		constant = "EDESTADDRREQ"
	case unix.EDOM:
		constant = "EDOM"
	case unix.EDOTDOT:
		constant = "EDOTDOT"
	case unix.EDQUOT:
		constant = "EDQUOT"
	case unix.EEXIST:
		constant = "EEXIST"
	case unix.EFAULT:
		constant = "EFAULT"
	case unix.EFBIG:
		constant = "EFBIG"
	case unix.EHOSTDOWN:
		constant = "EHOSTDOWN"
	case unix.EHOSTUNREACH:
		constant = "EHOSTUNREACH"
	case unix.EIDRM:
		constant = "EIDRM"
	case unix.EILSEQ:
		constant = "EILSEQ"
	case unix.EINPROGRESS:
		constant = "EINPROGRESS"
	case unix.EINTR:
		constant = "EINTR"
	case unix.EINVAL:
		constant = "EINVAL"
	case unix.EIO:
		constant = "EIO"
	case unix.EISCONN:
		constant = "EISCONN"
	case unix.EISDIR:
		constant = "EISDIR"
	case unix.EISNAM:
		constant = "EISNAM"
	case unix.EKEYEXPIRED:
		constant = "EKEYEXPIRED"
	case unix.EKEYREJECTED:
		constant = "EKEYREJECTED"
	case unix.EKEYREVOKED:
		constant = "EKEYREVOKED"
	case unix.EL2HLT:
		constant = "EL2HLT"
	case unix.EL2NSYNC:
		constant = "EL2NSYNC"
	case unix.EL3HLT:
		constant = "EL3HLT"
	case unix.EL3RST:
		constant = "EL3RST"
	case unix.ELIBACC:
		constant = "ELIBACC"
	case unix.ELIBBAD:
		constant = "ELIBBAD"
	case unix.ELIBEXEC:
		constant = "ELIBEXEC"
	case unix.ELIBMAX:
		constant = "ELIBMAX"
	case unix.ELIBSCN:
		constant = "ELIBSCN"
	case unix.ELNRNG:
		constant = "ELNRNG"
	case unix.ELOOP:
		constant = "ELOOP"
	case unix.EMEDIUMTYPE:
		constant = "EMEDIUMTYPE"
	case unix.EMFILE:
		constant = "EMFILE"
	case unix.EMLINK:
		constant = "EMLINK"
	case unix.EMSGSIZE:
		constant = "EMSGSIZE"
	case unix.EMULTIHOP:
		constant = "EMULTIHOP"
	case unix.ENAMETOOLONG:
		constant = "ENAMETOOLONG"
	case unix.ENAVAIL:
		constant = "ENAVAIL"
	case unix.ENETDOWN:
		constant = "ENETDOWN"
	case unix.ENETRESET:
		constant = "ENETRESET"
	case unix.ENETUNREACH:
		constant = "ENETUNREACH"
	case unix.ENFILE:
		constant = "ENFILE"
	case unix.ENOANO:
		constant = "ENOANO"
	case unix.ENOBUFS:
		constant = "ENOBUFS"
	case unix.ENOCSI:
		constant = "ENOCSI"
	case unix.ENODATA:
		constant = "ENODATA"
	case unix.ENODEV:
		constant = "ENODEV"
	case unix.ENOENT:
		constant = "ENOENT"
	case unix.ENOEXEC:
		constant = "ENOEXEC"
	case unix.ENOKEY:
		constant = "ENOKEY"
	case unix.ENOLCK:
		constant = "ENOLCK"
	case unix.ENOLINK:
		constant = "ENOLINK"
	case unix.ENOMEDIUM:
		constant = "ENOMEDIUM"
	case unix.ENOMEM:
		constant = "ENOMEM"
	case unix.ENOMSG:
		constant = "ENOMSG"
	case unix.ENONET:
		constant = "ENONET"
	case unix.ENOPKG:
		constant = "ENOPKG"
	case unix.ENOPROTOOPT:
		constant = "ENOPROTOOPT"
	case unix.ENOSPC:
		constant = "ENOSPC"
	case unix.ENOSR:
		constant = "ENOSR"
	case unix.ENOSTR:
		constant = "ENOSTR"
	case unix.ENOSYS:
		constant = "ENOSYS"
	case unix.ENOTBLK:
		constant = "ENOTBLK"
	case unix.ENOTCONN:
		constant = "ENOTCONN"
	case unix.ENOTDIR:
		constant = "ENOTDIR"
	case unix.ENOTEMPTY:
		constant = "ENOTEMPTY"
	case unix.ENOTNAM:
		constant = "ENOTNAM"
	case unix.ENOTRECOVERABLE:
		constant = "ENOTRECOVERABLE"
	case unix.ENOTSOCK:
		constant = "ENOTSOCK"
	case unix.ENOTSUP:
		constant = "ENOTSUP"
	case unix.ENOTTY:
		constant = "ENOTTY"
	case unix.ENOTUNIQ:
		constant = "ENOTUNIQ"
	case unix.ENXIO:
		constant = "ENXIO"
	case unix.EOVERFLOW:
		constant = "EOVERFLOW"
	case unix.EOWNERDEAD:
		constant = "EOWNERDEAD"
	case unix.EPERM:
		constant = "EPERM"
	case unix.EPFNOSUPPORT:
		constant = "EPFNOSUPPORT"
	case unix.EPIPE:
		constant = "EPIPE"
	case unix.EPROTO:
		constant = "EPROTO"
	case unix.EPROTONOSUPPORT:
		constant = "EPROTONOSUPPORT"
	case unix.EPROTOTYPE:
		constant = "EPROTOTYPE"
	case unix.ERANGE:
		constant = "ERANGE"
	case unix.EREMCHG:
		constant = "EREMCHG"
	case unix.EREMOTE:
		constant = "EREMOTE"
	case unix.EREMOTEIO:
		constant = "EREMOTEIO"
	case unix.ERESTART:
		constant = "ERESTART"
	case unix.ERFKILL:
		constant = "ERFKILL"
	case unix.EROFS:
		constant = "EROFS"
	case unix.ESHUTDOWN:
		constant = "ESHUTDOWN"
	case unix.ESOCKTNOSUPPORT:
		constant = "ESOCKTNOSUPPORT"
	case unix.ESPIPE:
		constant = "ESPIPE"
	case unix.ESRCH:
		constant = "ESRCH"
	case unix.ESRMNT:
		constant = "ESRMNT"
	case unix.ESTALE:
		constant = "ESTALE"
	case unix.ESTRPIPE:
		constant = "ESTRPIPE"
	case unix.ETIME:
		constant = "ETIME"
	case unix.ETIMEDOUT:
		constant = "ETIMEDOUT"
	case unix.ETOOMANYREFS:
		constant = "ETOOMANYREFS"
	case unix.ETXTBSY:
		constant = "ETXTBSY"
	case unix.EUCLEAN:
		constant = "EUCLEAN"
	case unix.EUNATCH:
		constant = "EUNATCH"
	case unix.EUSERS:
		constant = "EUSERS"
	case unix.EXDEV:
		constant = "EXDEV"
	case unix.EXFULL:
		constant = "EXFULL"
	}

	msg := syscall.Errno(code).Error()
	p.PrintColour(ColourRed, "%d %s (%s)", arg.Int()+1, constant, msg)
}

package annotation

/*
#include <asm/signal.h>
#include <asm/siginfo.h>

*/
import "C"
import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func AnnotateSigProcMaskFlags(arg Arg, _ int) {
	var joins []string
	if arg.Raw()&C.SIG_BLOCK != 0 {
		joins = append(joins, "SIG_BLOCK")
	}
	if arg.Raw()&C.SIG_UNBLOCK != 0 {
		joins = append(joins, "SIG_UNBLOCK")
	}
	if arg.Raw()&C.SIG_SETMASK != 0 {
		joins = append(joins, "SIG_SETMASK")
	}
	arg.SetAnnotation(strings.Join(joins, "|"), len(joins) > 0)
}

var signals = map[int]string{
	C.SIGHUP:    "SIGHUP",
	C.SIGINT:    "SIGINT",
	C.SIGQUIT:   "SIGQUIT",
	C.SIGILL:    "SIGILL",
	C.SIGTRAP:   "SIGTRAP",
	C.SIGABRT:   "SIGABRT",
	C.SIGBUS:    "SIGBUS",
	C.SIGFPE:    "SIGFPE",
	C.SIGKILL:   "SIGKILL",
	C.SIGUSR1:   "SIGUSR1",
	C.SIGSEGV:   "SIGSEGV",
	C.SIGUSR2:   "SIGUSR2",
	C.SIGPIPE:   "SIGPIPE",
	C.SIGALRM:   "SIGALRM",
	C.SIGTERM:   "SIGTERM",
	C.SIGSTKFLT: "SIGSTKFLT",
	C.SIGCHLD:   "SIGCHLD",
	C.SIGCONT:   "SIGCONT",
	C.SIGSTOP:   "SIGSTOP",
	C.SIGTSTP:   "SIGTSTP",
	C.SIGTTIN:   "SIGTTIN",
	C.SIGTTOU:   "SIGTTOU",
	C.SIGURG:    "SIGURG",
	C.SIGXCPU:   "SIGXCPU",
	C.SIGXFSZ:   "SIGXFSZ",
	C.SIGVTALRM: "SIGVTALRM",
	C.SIGPROF:   "SIGPROF",
	C.SIGWINCH:  "SIGWINCH",
	C.SIGIO:     "SIGIO",
	C.SIGPWR:    "SIGPWR",
	C.SIGSYS:    "SIGSYS",
}

func AnnotateSignal(arg Arg, _ int) {
	arg.SetAnnotation(SignalToString(int(arg.Raw())), true)
}

func SignalToString(signal int) string {
	var extra string
	if signal&0x80 != 0 {
		extra = "|0x80"
		signal &= ^0x80
	}
	if name, ok := signals[signal]; ok {
		return name + extra
	}
	return fmt.Sprintf("0x%x%s", signal, extra)
}

var signalCodes = map[syscall.Signal]map[int]string{
	0: {
		C.SI_USER:    "SI_USER",
		C.SI_KERNEL:  "SI_KERNEL",
		C.SI_QUEUE:   "SI_QUEUE",
		C.SI_TIMER:   "SI_TIMER",
		C.SI_MESGQ:   "SI_MESGQ",
		C.SI_ASYNCIO: "SI_ASYNCIO",
		C.SI_SIGIO:   "SI_SIGIO",
		C.SI_TKILL:   "SI_TKILL",
	},
	syscall.SIGILL: {
		C.ILL_ILLOPC: "ILL_ILLOPC",
		C.ILL_ILLOPN: "ILL_ILLOPN",
		C.ILL_ILLADR: "ILL_ILLADR",
		C.ILL_ILLTRP: "ILL_ILLTRP",
		C.ILL_PRVOPC: "ILL_PRVOPC",
		C.ILL_PRVREG: "ILL_PRVREG",
		C.ILL_COPROC: "ILL_COPROC",
		C.ILL_BADSTK: "ILL_BADSTK",
	},
	syscall.SIGFPE: {
		C.FPE_INTDIV: "FPE_INTDIV",
		C.FPE_INTOVF: "FPE_INTOVF",
		C.FPE_FLTDIV: "FPE_FLTDIV",
		C.FPE_FLTOVF: "FPE_FLTOVF",
		C.FPE_FLTUND: "FPE_FLTUND",
		C.FPE_FLTRES: "FPE_FLTRES",
		C.FPE_FLTINV: "FPE_FLTINV",
		C.FPE_FLTSUB: "FPE_FLTSUB",
	},
	syscall.SIGSEGV: {
		C.SEGV_MAPERR: "SEGV_MAPERR",
		C.SEGV_ACCERR: "SEGV_ACCERR",
		C.SEGV_BNDERR: "SEGV_BNDERR",
		C.SEGV_PKUERR: "SEGV_PKUERR",
	},
	syscall.SIGBUS: {
		C.BUS_ADRALN:    "BUS_ADRALN",
		C.BUS_ADRERR:    "BUS_ADRERR",
		C.BUS_OBJERR:    "BUS_OBJERR",
		C.BUS_MCEERR_AR: "BUS_MCEERR_AR",
		C.BUS_MCEERR_AO: "BUS_MCEERR_AO",
	},
	syscall.SIGTRAP: {
		C.TRAP_BRKPT:  "TRAP_BRKPT",
		C.TRAP_TRACE:  "TRAP_TRACE",
		C.TRAP_BRANCH: "TRAP_BRANCH",
		C.TRAP_HWBKPT: "TRAP_HWBKPT",
	},
	syscall.SIGCHLD: {
		C.CLD_EXITED:    "CLD_EXITED",
		C.CLD_KILLED:    "CLD_KILLED",
		C.CLD_DUMPED:    "CLD_DUMPED",
		C.CLD_TRAPPED:   "CLD_TRAPPED",
		C.CLD_STOPPED:   "CLD_STOPPED",
		C.CLD_CONTINUED: "CLD_CONTINUED",
	},
	syscall.SIGPOLL: {
		C.POLL_IN:  "POLL_IN",
		C.POLL_OUT: "POLL_OUT",
		C.POLL_MSG: "POLL_MSG",
		C.POLL_ERR: "POLL_ERR",
		C.POLL_PRI: "POLL_PRI",
		C.POLL_HUP: "POLL_HUP",
	},
	syscall.SIGSYS: {
		C.SYS_SECCOMP: "SYS_SECCOMP",
	},
}

func SignalCodeToString(signal syscall.Signal, code int32) string {
	if codes, ok := signalCodes[signal]; ok {
		if str, ok := codes[int(code)]; ok {
			return str
		}
	}
	if str, ok := signalCodes[0][int(code)]; ok {
		return str
	}
	return fmt.Sprintf("%d", code)
}

var signalFDFlags = map[int]string{
	unix.SFD_NONBLOCK: "SFD_NONBLOCK",
	unix.SFD_CLOEXEC:  "SFD_CLOEXEC",
}

func AnnotateSignalFdFlags(arg Arg, _ int) {
	if str, ok := signalFDFlags[int(arg.Raw())]; ok {
		arg.SetAnnotation(str, true)
	}
}

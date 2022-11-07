package printer

import (
	"syscall"

	"github.com/liamg/grace/tracer"
	"github.com/liamg/grace/tracer/annotation"
)

func (p *Printer) PrintSignal(signal *tracer.SigInfo) {
	p.PrefixEvent()
	p.PrintColour(ColourMagenta, "--> ")
	p.PrintColour(
		ColourCyan,
		"SIGNAL: %s (code=%s, pid=%d, uid=%d)",
		annotation.SignalToString(int(signal.Signo)),
		annotation.SignalCodeToString(syscall.Signal(signal.Signo), signal.Code),
		signal.Pid,
		signal.Uid,
	)
	p.PrintColour(ColourMagenta, " <--\n")
	if p.multiline {
		p.Print("\n")
	}
}

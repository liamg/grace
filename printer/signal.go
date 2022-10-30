package printer

import (
	"fmt"
	"syscall"
)

func (p *Printer) PrintSignal(signal syscall.Signal) {
	p.PrintColour(ColourYellow, "--> SIGNAL: %s <--\n", signalToString(signal))
}

func signalToString(signal syscall.Signal) string {
	switch signal {
	// TODO: more signals
	case syscall.SIGURG:
		return "SIGURG"
	default:
		return fmt.Sprintf("signal %d", signal)
	}
}

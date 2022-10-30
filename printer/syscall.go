package printer

import "github.com/liamg/grace/tracer"

func (p *Printer) PrintSyscallEnter(syscall *tracer.Syscall) {
	p.colourIndex = 0
	p.argProgress = 0
	if syscall.Unknown() {
		p.PrintColour(ColourRed, syscall.Name())
	} else {
		p.PrintColour(ColourDefault, syscall.Name())
	}
	p.printRemainingArgs(syscall, false)
	p.inSyscall = true
}

func (p *Printer) PrintSyscallExit(syscall *tracer.Syscall) {
	p.printRemainingArgs(syscall, true)
	p.PrintDim(" = ")
	ret := syscall.Return()
	p.PrintArgValue(&ret, ColourWhite, true, 0, 0)
	p.Print("\n")
	if p.extraNewLine {
		p.Print("\n")
	}
	p.inSyscall = false
}

func (p *Printer) printRemainingArgs(syscall *tracer.Syscall, exit bool) {
	if !exit {
		p.PrintDim("(")
	}
	remaining := syscall.Args()[p.argProgress:]
	for i, arg := range remaining {
		if !arg.Known() {
			break
		}
		if p.argProgress == 0 && p.multiline {
			p.Print("\n")
		}
		p.PrintArg(arg, exit)
		if i < len(remaining)-1 {
			p.PrintDim(", ")
		}
		p.argProgress++
		if p.multiline {
			p.Print("\n")
		}
	}

	if (exit && len(remaining) > 0) || (!exit && p.argProgress == len(syscall.Args())) {
		p.PrintDim(")")
	}

}

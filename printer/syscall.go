package printer

import (
	"fmt"

	"github.com/liamg/grace/tracer"
)

func (p *Printer) PrintSyscallEnter(syscall *tracer.Syscall) {
	p.printSyscallEnter(syscall, false)
}

func (p *Printer) printSyscallEnter(syscall *tracer.Syscall, overrideFilter bool) {

	if !overrideFilter {
		if p.filter != nil {
			if !p.filter.Match(syscall, false) {
				p.lastEntryMatchedFilter = false
				return
			}
		}
		p.lastEntryMatchedFilter = true
	}

	p.PrefixEvent()

	p.colourIndex = 0
	p.argProgress = 0

	if p.showNumbers {
		p.PrintDim("%4s", fmt.Sprintf("%d ", syscall.Number()))
	}

	if syscall.Unknown() {
		p.PrintColour(ColourRed, syscall.Name())
	} else {
		p.PrintColour(ColourDefault, syscall.Name())
	}
	p.printRemainingArgs(syscall, false)
	p.inSyscall = true
}

func (p *Printer) PrintSyscallExit(syscall *tracer.Syscall) {

	if p.filter != nil {
		if !p.lastEntryMatchedFilter && !p.filter.Match(syscall, true) {
			return
		}
	}

	if !p.lastEntryMatchedFilter {
		p.printSyscallEnter(syscall, true)
	}

	p.printRemainingArgs(syscall, true)
	p.PrintDim(" = ")
	ret := syscall.Return()
	p.PrintArgValue(&ret, ColourGreen, true, 0, 0)
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
	var remaining []tracer.Arg
	if p.argProgress < len(syscall.Args()) {
		remaining = syscall.Args()[p.argProgress:]
		for i, arg := range remaining {
			if !arg.Known() {
				break
			}
			if p.argProgress == 0 && p.multiline {
				p.Print("\n")
			}
			p.PrintArg(arg, exit)
			if i < len(remaining)-1 || !syscall.Complete() {
				p.PrintDim(", ")
			}
			p.argProgress++
			if p.multiline {
				p.Print("\n")
			}
		}
	}

	if ((exit && len(remaining) > 0) || (!exit && p.argProgress == len(syscall.Args()))) && syscall.Complete() {
		p.PrintDim(")")
	}

}

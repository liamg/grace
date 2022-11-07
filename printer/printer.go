package printer

import (
	"fmt"
	"io"
	"time"

	"github.com/liamg/grace/tracer"
)

type Printer struct {
	w                      io.Writer
	useColours             bool
	maxStringLen           int
	hexDumpLongStrings     bool
	maxHexDumpLen          int
	maxObjectProperties    int
	colourIndex            int
	argProgress            int
	extraNewLine           bool
	multiline              bool
	inSyscall              bool
	filter                 Filter
	lastEntryMatchedFilter bool
	relativeTimestamps     bool
	absoluteTimestamps     bool
	startTime              time.Time
	showNumbers            bool
	rawOutput              bool
}

type Filter interface {
	Match(syscall *tracer.Syscall, exit bool) bool
}

func New(w io.Writer) *Printer {
	return &Printer{
		w:                   w,
		useColours:          true,
		maxStringLen:        32,
		hexDumpLongStrings:  true,
		maxHexDumpLen:       4096,
		maxObjectProperties: 2,
		startTime:           time.Now(),
	}
}

const indentSize = 4

func (p *Printer) SetUseColours(useColours bool) {
	p.useColours = useColours
}

func (p *Printer) SetMaxStringLen(maxStringLen int) {
	p.maxStringLen = maxStringLen
}

func (p *Printer) SetHexDumpLongStrings(hexDumpLongStrings bool) {
	p.hexDumpLongStrings = hexDumpLongStrings
}

func (p *Printer) SetMaxHexDumpLen(maxHexDumpLen int) {
	p.maxHexDumpLen = maxHexDumpLen
}

func (p *Printer) SetMaxObjectProperties(maxObjectProperties int) {
	p.maxObjectProperties = maxObjectProperties
}

func (p *Printer) SetExtraNewLine(extraNewLine bool) {
	p.extraNewLine = extraNewLine
}

func (p *Printer) SetShowAbsoluteTimestamps(timestamps bool) {
	p.absoluteTimestamps = timestamps
}

func (p *Printer) SetShowRelativeTimestamps(timestamps bool) {
	p.relativeTimestamps = timestamps
}

func (p *Printer) SetMultiLine(multiline bool) {
	p.multiline = multiline
}

func (p *Printer) SetFilter(filter Filter) {
	p.filter = filter
}

func (p *Printer) SetShowSyscallNumber(number bool) {
	p.showNumbers = number
}

func (p *Printer) SetRawOutput(output bool) {
	p.rawOutput = output
}

func (p *Printer) PrefixEvent() {
	if p.relativeTimestamps {
		p.PrintDim("%12s ", time.Since(p.startTime))
	}
	if p.absoluteTimestamps {
		p.PrintDim("%18s ", time.Now().Format("15:04:05.999999999"))
	}
}

func (p *Printer) Print(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.w, format, args...)
}

func (p *Printer) PrintDim(format string, args ...interface{}) {
	p.PrintColour(2, format, args...)
}

func (p *Printer) PrintColour(colour Colour, format string, args ...interface{}) {
	if p.useColours {
		p.Print("\x1b[%dm", colour)
	}
	p.Print(format, args...)
	if p.useColours {
		p.Print("\x1b[0m")
	}
}

func (p *Printer) PrintProcessExit(i int) {
	colour := ColourGreen
	if i != 0 {
		colour = ColourRed
	}
	if p.inSyscall {
		p.PrintDim(" = ?\n")
	}
	if p.multiline {
		p.Print("\n")
	}
	p.PrintColour(colour, "Process exited with status %d\n", i)
}

func (p *Printer) PrintAttach(pid int) {
	p.PrintColour(ColourYellow, "Attached to process %d\n", pid)
	if p.multiline {
		p.Print("\n")
	}
}

func (p *Printer) PrintDetach(pid int) {
	p.PrintColour(ColourYellow, "Detached from process %d\n", pid)
	if p.multiline {
		p.Print("\n")
	}
}

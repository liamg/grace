package printer

import (
	"fmt"
	"io"
)

type Printer struct {
	w                   io.Writer
	useColours          bool
	maxStringLen        int
	hexDumpLongStrings  bool
	maxHexDumpLen       int
	maxObjectProperties int
	colourIndex         int
	argProgress         int
	extraNewLine        bool
	multiline           bool
}

func New(w io.Writer) *Printer {
	return &Printer{
		w:                   w,
		useColours:          true,
		maxStringLen:        32,
		hexDumpLongStrings:  true,
		maxHexDumpLen:       4096,
		maxObjectProperties: 2,
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

func (p *Printer) SetMultiLine(multiline bool) {
	p.multiline = multiline
}

func (p *Printer) PrintProcessExit(i int) {
	colour := ColourGreen
	if i != 0 {
		colour = ColourRed
	}
	if p.multiline {
		p.Print("\n")
	}
	p.PrintColour(colour, "\nProcess exited with status %d\n", i)
}

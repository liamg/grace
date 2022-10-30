package printer

import "strings"

const dumpWidth = 16

func (p *Printer) HexDump(addr uintptr, data []byte, indent int) {

	var truncatedFrom uintptr
	if p.maxHexDumpLen > 0 && len(data) > p.maxHexDumpLen {
		truncatedFrom = uintptr(len(data))
		data = data[:p.maxHexDumpLen]
	}

	startAddr := addr - (addr % dumpWidth)
	endAddr := addr + uintptr(len(data))
	if endAddr%dumpWidth > 0 {
		endAddr += dumpWidth - (endAddr % dumpWidth)
	}
	p.Print(strings.Repeat(" ", indent))
	p.Print("(see below hexdump)")
	p.NewLine(indent)
	p.Print("                  ")
	for i := 0; i < dumpWidth; i++ {
		p.PrintDim("%02x ", i)
	}
	for i := startAddr; i < endAddr; i += dumpWidth {
		p.NewLine(indent)
		p.PrintDim("%16x: ", i)
		for j := 0; j < dumpWidth; j++ {
			local := (i + uintptr(j)) - addr
			if i+uintptr(j) < addr || local >= uintptr(len(data)) {
				p.PrintDim(".. ")
			} else {
				p.PrintColour(ColourRed, "%02x ", data[local])
			}
		}
		for j := 0; j < dumpWidth; j++ {
			local := (i + uintptr(j)) - addr
			if i+uintptr(j) < addr || local >= uintptr(len(data)) {
				p.PrintDim(".")
			} else {
				c := data[local]
				if c < 32 || c > 126 {
					c = '.'
				}
				p.PrintColour(ColourBlue, "%c", c)
			}
		}
	}
	if truncatedFrom > 0 {
		p.PrintDim("\n... (truncated from %d bytes -> %d bytes) ...", truncatedFrom, p.maxHexDumpLen)
	}
}

package printer

import "github.com/liamg/grace/tracer"

func (p *Printer) printArray(arr []tracer.Arg, colour Colour, exit bool, count int, indent int) int {
	prevIndent := indent
	indent += indentSize
	p.PrintColour(colour, "[")

	for i, prop := range arr {
		if p.multiline {
			p.NewLine(indent)
		}
		colour := colours[i%len(colours)]
		p.PrintDim("%d", i)
		p.PrintDim(": ")
		count += p.PrintArgValue(&prop, colour, exit, count, indent+indentSize)
		if i < len(arr)-1 {
			p.PrintDim(", ")
		}
		if p.maxObjectProperties > 0 && count >= p.maxObjectProperties && count < len(arr) {
			p.PrintDim("...")
			break
		}
		count++
	}
	if p.multiline && len(arr) > 0 {
		p.NewLine(prevIndent)
	}
	p.PrintColour(colour, "]")
	return count
}

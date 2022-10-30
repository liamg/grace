package printer

import (
	"github.com/liamg/grace/tracer"
)

func (p *Printer) printObject(obj *tracer.Object, colour Colour, exit bool, count int, indent int) int {
	if obj == nil {
		p.PrintDim("NULL")
		return count
	}
	prevIndent := indent
	indent += indentSize
	p.PrintColour(colour, "%s{", obj.Name)

	for i, prop := range obj.Properties {
		if p.multiline {
			p.NewLine(indent)
		}
		colour := colours[i%len(colours)]
		p.PrintDim("%s", prop.Name())
		p.PrintDim(": ")
		count += p.PrintArgValue(&prop, colour, exit, count, indent+indentSize)
		if i < len(obj.Properties)-1 {
			p.PrintDim(", ")
		}
		if p.maxObjectProperties > 0 && count >= p.maxObjectProperties {
			p.PrintDim("...")
			break
		}
		count++
	}
	if p.multiline {
		p.NewLine(prevIndent)
	}
	p.PrintColour(colour, "}")
	return count
}

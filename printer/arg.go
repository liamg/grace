package printer

import (
	"strings"

	"github.com/liamg/grace/tracer"
)

func (p *Printer) PrintArg(arg tracer.Arg, exit bool) {

	var indent int
	if p.multiline {
		indent = indentSize
	}

	if p.multiline {
		p.Print(strings.Repeat(" ", indent))
	}

	if name := arg.Name(); name != "" {
		p.PrintDim("%s: ", name)
	}

	p.PrintArgValue(&arg, p.nextColour(), exit, 0, indent)
}

func (p *Printer) NewLine(indent int) {
	p.Print("\n" + strings.Repeat(" ", indent))
}

func (p *Printer) PrintArgValue(arg *tracer.Arg, colour Colour, exit bool, propCount int, indent int) int {

	if arg.ReplaceValueWithAnnotation() {
		p.PrintColour(colour, "%s", arg.Annotation())
		return propCount
	}

	switch arg.Type() {
	case tracer.ArgTypeData:
		data := arg.Data()
		if p.maxStringLen > 0 && len(data) > p.maxStringLen {
			if exit && p.hexDumpLongStrings {
				p.HexDump(arg.Raw(), arg.Data(), indent)
				return propCount
			}
			data = append(data[:p.maxStringLen], []byte("...")...)
		}
		p.PrintColour(colour, "%q", string(data))
		//p.PrintDim(" @ 0x%x", arg.Raw())
	case tracer.ArgTypeInt, tracer.ArgTypeLong, tracer.ArgTypeUnsignedInt, tracer.ArgTypeUnsignedLong, tracer.ArgTypeUnknown:
		p.PrintColour(colour, "%d", arg.Int())
	case tracer.ArgTypeErrorCode:
		p.printError(colour, arg)
	case tracer.ArgTypeAddress:
		p.PrintColour(colour, "0x%x", arg.Raw())
	case tracer.ArgTypeObject, tracer.ArgTypeStat, tracer.ArgTypeSigAction:
		propCount += p.printObject(arg.Object(), colour, exit, propCount, indent)
	case tracer.ArgTypeIntArray, tracer.ArgTypePollFdArray, tracer.ArgTypeIovecArray:
		propCount += p.printArray(arg.Array(), colour, exit, propCount, indent)
	default:
		// TODO: error here
		p.PrintColour(colour, "%d (?)", arg.Raw())
	}

	if annotation := arg.Annotation(); annotation != "" {
		p.PrintDim(" -> %s", annotation)
	}

	return propCount
}

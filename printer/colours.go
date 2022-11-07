package printer

type Colour int

const (
	ColourRed Colour = iota + 31
	ColourGreen
	ColourYellow
	ColourBlue
	ColourMagenta
	ColourCyan
	ColourWhite
)

var ColourDefault Colour = 0

var colours = []Colour{
	ColourBlue,
	ColourYellow,
	ColourGreen,
}

func (p *Printer) currentColour() Colour {
	return colours[p.colourIndex%len(colours)]
}

func (p *Printer) nextColour() Colour {
	colour := colours[p.colourIndex%len(colours)]
	p.colourIndex++
	return colour
}

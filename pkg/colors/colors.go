package colors

import (
	"github.com/lucasb-eyer/go-colorful"
)

const MaxColorModes = 3

type gradientTableEntry struct {
	Col colorful.Color
	Pos float64
}

type GradientTable struct {
	table []gradientTableEntry
	Mode  int
}

// Parses hex strings into colors
func ParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("parseHex: " + err.Error())
	}
	return c
}

func (gt *GradientTable) AddToTable(colorString string, pos float64) {
	gt.table = append(gt.table, gradientTableEntry{ParseHex(colorString), pos})
}

func (gt *GradientTable) Randomise() {
	startEnd := colorful.FastHappyColor()
	gt.table = nil
	gt.table = append(gt.table, gradientTableEntry{startEnd, 0.00})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.25})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.50})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.75})
	gt.table = append(gt.table, gradientTableEntry{startEnd, 1.00})
}

// Get a blended color
func (gt *GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt.table)-1; i++ {
		c1 := gt.table[i]
		c2 := gt.table[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			switch gt.Mode {
			case 0:
				return c1.Col.BlendRgb(c2.Col, t).Clamped()
			case 1:
				return c1.Col.BlendHsv(c2.Col, t).Clamped()
			case 2:
				return c1.Col.BlendHcl(c2.Col, t).Clamped()
			}
		}
	}

	// This deals with wrapping around past end of last entry
	first := gt.table[0]
	last := gt.table[len(gt.table)-1]
	wrappedT := (t - last.Pos) / (1.0 - last.Pos)
	return last.Col.BlendRgb(first.Col, wrappedT).Clamped()

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	// return gt.table[len(gt.table)-1].Col
}

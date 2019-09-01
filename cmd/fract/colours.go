package main

import "github.com/lucasb-eyer/go-colorful"

const maxColorModes = 3

type gradientTableEntry struct {
	Col colorful.Color
	Pos float64
}

type gradientTable struct {
	table []gradientTableEntry
	mode  int
}

// Parses hex strings into colors
func parseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("parseHex: " + err.Error())
	}
	return c
}

func (gt *gradientTable) addToTable(colorString string, pos float64) {
	gt.table = append(gt.table, gradientTableEntry{parseHex(colorString), pos})
}

func (gt *gradientTable) randomise() {
	gt.table = nil
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.0})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.333})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 0.666})
	gt.table = append(gt.table, gradientTableEntry{colorful.FastHappyColor(), 1.0})
}

// Get a blended color
func (gt *gradientTable) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt.table)-1; i++ {
		c1 := gt.table[i]
		c2 := gt.table[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			switch gt.mode {
			case 0:
				return c1.Col.BlendRgb(c2.Col, t).Clamped()
			case 1:
				return c1.Col.BlendHsv(c2.Col, t).Clamped()
			case 2:
				return c1.Col.BlendHcl(c2.Col, t).Clamped()
			}
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt.table[len(gt.table)-1].Col
}

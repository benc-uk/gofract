package main

import "github.com/lucasb-eyer/go-colorful"

type gradientTable []struct {
	Col colorful.Color
	Pos float64
}

func parseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("parseHex: " + err.Error())
	}
	return c
}

// blah
func (gt gradientTable) getInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

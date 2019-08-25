package main

import (
	"math"
)

var escape = 4.0
var log2 = math.Log(2)
var escapeOffset = math.Log(math.Log(escape))

func mandlebrot(r, i float64) int {
	var x = 0.0
	var y = 0.0

	var iter = 0

	for x*x+y*y <= escape && iter < maxIter {
		xtemp := x*x - y*y + r
		y = 2*x*y + i
		x = xtemp
		iter = iter + 1
	}
	if iter == maxIter {
		return maxIter
	}

	m := (float64(iter) - escapeOffset) / log2 
	return int(m)
}

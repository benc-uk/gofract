package main

import (
	"math"
	"math/cmplx"
)

var (
	escape = 2.0
	log2   = math.Log(escape)
)

func mandlebrot(a complex128, f Fractal) int {
	var z complex128 // zero
	var iter = 0

	for cmplx.Abs(z) < escape && iter <= f.maxIter {
		z = z * z + a
		iter++
	}

	if iter >= f.maxIter {
		return f.maxIter
	}

	mu := float64(iter) + 2.0 - math.Log(math.Log(cmplx.Abs(z)))/log2
	return int(mu)
}

func julia(a complex128, f Fractal) int {
	z := a
	var iter = 0

	for cmplx.Abs(z) < escape && iter <= f.maxIter {
		z = z * z + f.c
		iter++
	}

	if iter >= f.maxIter {
		return f.maxIter
	}

	mu := float64(iter) + 2.0 - math.Log(math.Log(cmplx.Abs(z)))/log2
	return int(mu)
}

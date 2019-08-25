package main

import (
	//"fmt"
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
)

var mainCanvas *canvas.Raster
var gradient gradientTable
var bg = color.RGBA{0,0,0,255}

var maxIter = 50
var magFactor = 40.0
var centerX = 3.0
var centerY = 0.0

var xWidth = 3.0
var yHeight = 2.0

func main() {
	gradient = gradientTable{
		{parseHex("#9e0142"), 0.0},
		{parseHex("#d53e4f"), 0.1},
		{parseHex("#f46d43"), 0.2},
		{parseHex("#fdae61"), 0.3},
		{parseHex("#fee090"), 0.4},
		{parseHex("#ffffbf"), 0.5},
		{parseHex("#e6f598"), 0.6},
		{parseHex("#abdda4"), 0.7},
		{parseHex("#66c2a5"), 0.8},
		{parseHex("#3288bd"), 0.9},
		{parseHex("#5e4fa2"), 1.0},
	}

	fmt.Println("### Starting...")
	app := app.New()
	window := app.NewWindow("GoFract")
	window.SetPadded(false)
	mainCanvas = canvas.NewRasterWithPixels(drawFractal)
	window.SetContent(mainCanvas)
	// window.SetFixedSize(true)
	window.Show()
	window.Resize(fyne.Size{1200, 800})
	mainCanvas.Resize(fyne.Size{1200, 800})

	window.ShowAndRun()
}

func drawFractal(x, y, w, h int) color.Color {
	RE_START := -2.0
	RE_END := 1.0
	IM_START := -1.0
	IM_END := 1.0

	if w == 0 || h == 0 {
		return bg
	}
	wf := float64(w)
	hf := float64(h)
	// ar := 1.0
	// if(w > h) {
	// 	ar = wf/hf
	// } else {
	// 	ar = hf/wf
	// }

	r := RE_START + ((float64(x) / wf) * (RE_END - RE_START))
	i := IM_START + ((float64(y) / hf) * (IM_END - IM_START))
	
	//fmt.Println(w)
	iter := mandlebrot(float64(r), float64(i))

	if iter == maxIter {
		//return color.Black
		return bg
	}

	//return color.RGBA{255,10,20,255}
	c := float64(iter) / float64(maxIter)
	return gradient.getInterpolatedColorFor(c)
}

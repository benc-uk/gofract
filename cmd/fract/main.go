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

var maxIter = 80
var magFactor = 3.0
var rOffset = -2.0
var iOffset = -1.0
var imgWidth = 2000
var imgHeight = 1500

func main() {
	gradient = gradientTable{
		{parseHex("#090cb5"), 0.0},
		// {parseHex("#d53e4f"), 0.1},
		{parseHex("#7627ab"), 0.15},
		// {parseHex("#fdae61"), 0.3},
		// {parseHex("#fee090"), 0.4},
		{parseHex("#f56320"), 0.4},
		// {parseHex("#e6f598"), 0.6},
		// {parseHex("#abdda4"), 0.7},
		// {parseHex("#66c2a5"), 0.8},
		// {parseHex("#3288bd"), 0.9},
		{parseHex("#ffff00"), 1.0},
	}

	fmt.Println("### Starting GoFract...")
	app := app.New()
	window := app.NewWindow("GoFract")
	window.SetPadded(false)
	mainCanvas = canvas.NewRasterWithPixels(drawFractal)
	window.SetContent(mainCanvas)
	window.Canvas().SetOnTypedKey(keyPressed)
	window.Show()
	fmt.Println("### WINDOW SIZE:", imgWidth, imgHeight)

	window.Resize(fyne.Size{ imgWidth, imgHeight })
	mainCanvas.Resize(fyne.Size{ imgWidth, imgHeight })

	window.ShowAndRun()
}

func drawFractal(x, y, w, h int) color.Color {
	wf := float64(w)
	hf := float64(h)

	rRatioAdjust := 1.0
	iRatioAdjust := 1.0
	if(w < h) {
		rRatioAdjust = wf / hf
	} else {
		iRatioAdjust = hf / wf
	}

	r := rOffset + ((float64(x) / wf) * magFactor * rRatioAdjust)
	i := iOffset + ((float64(y) / hf) * magFactor * iRatioAdjust)
	iter := mandlebrot(float64(r), float64(i))

	if iter == maxIter {
		return bg
	}

	scaledIter := float64(iter) / float64(maxIter)
	return gradient.getInterpolatedColorFor(scaledIter)
}

func keyPressed(ev *fyne.KeyEvent) {
	fmt.Println(ev.Name)
}
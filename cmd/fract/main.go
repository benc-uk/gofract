package main

import (
	"fmt"
	"image/color"
	// "golang.org/x/mobile/event/key"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
)

//var mainCanvas *canvas.Raster
var mainWidget fractWidget
var gradient gradientTable
var bg = color.RGBA{0,0,0,255}

var maxIter = 80
var magFactor = 3.0
var rOffset = -2.0
var iOffset = -1.0
var imgWidth = 1000
var imgHeight = 600

type fractWidget struct {
	canvas *canvas.Raster
	fractType string
	win fyne.Window
}

//
func main() {
	gradient = gradientTable{
		{parseHex("#090cb5"), 0.0},
		{parseHex("#7627ab"), 0.1},
		{parseHex("#f56320"), 0.3},
		{parseHex("#ffff00"), 1.0},
	}

	fmt.Println("### Starting GoFract...")
	app := app.New()
	window := app.NewWindow("GoFract")
	window.SetPadded(false)
	//mainCanvas = canvas.NewRasterWithPixels(drawFractal)

	mainWidget = fractWidget{
		canvas: canvas.NewRasterWithPixels(drawFractal),
		fractType: "mandelbrot",
		win: window,
	}

	window.SetContent(mainWidget.canvas)
	window.Canvas().SetOnTypedKey(keyEvent)
	window.Canvas().SetOnTypedRune(runeEvent)
	fmt.Println("### WINDOW SIZE:", imgWidth, imgHeight)

	window.Resize(fyne.Size{ imgWidth, imgHeight })
	mainWidget.canvas.Resize(fyne.Size{ imgWidth, imgHeight })

	window.ShowAndRun()
}

//
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

func keyEvent(ev *fyne.KeyEvent) {
	f := 0.1
	if ev.Name == "Up" {
		iOffset -= (magFactor * f)
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}
	if ev.Name == "Down" {
		iOffset += (magFactor * f)
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}
	if ev.Name == "Left" {
		rOffset -= (magFactor * f)
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}
	if ev.Name == "Right" {
		rOffset += (magFactor * f)
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}
}

func runeEvent(r int32) {
	if r == 61 {
		magFactor *= 0.9
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}

	if r == 45 {
		magFactor *= 1.1
		mainWidget.win.Canvas().Refresh(mainWidget.canvas)
	}
}


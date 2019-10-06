package main

import (
	"fmt"
	"image/color"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	exampleFonts "github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var helpText = []string{
	"HELP - Controls and keys:",
	" - Click anywhere to recenter & move",
	" - Mouse wheel up/down to zoom in and out",
	" - press 'h' to show this help text",
	" - press 'd' to show debug info",
	" - press 's' to save current view as PNG",
	" - press 'b' to change colour blend mode (RGB, HSV, HCL)",
	" - press 'm' to randomize the colour pallette",
	" - press 'z' and 'x' to decrease/increase max iterations",
	" - press 'r' to reload/reset fractal (will reload YAML from disk)",
	" - cursor keys to modify real/imaginary values of C when in Julia set mode",
}

//
// Render help text
//
func drawHelpOverlay(dst *ebiten.Image) {
	var fontSize = int(float64(f.ImgWidth) / 40.0)
	var ttf, _ = truetype.Parse(exampleFonts.MPlus1pRegular_ttf)

	var font = truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	for ix, textString := range helpText {
		text.Draw(dst, textString, font, 20+2, ((ix+3)*(fontSize+6))+2, color.Black)
		text.Draw(dst, textString, font, 20, (ix+3)*(fontSize+6), color.White)
	}
}

//
// Welcome message, only shown for the first 120 ticks
//
var welcomeCounter = 120

func drawWelcomeOverlay(dst *ebiten.Image) {
	var fontSize = int(float64(f.ImgWidth) / 60.0)
	var ttf, _ = truetype.Parse(exampleFonts.ArcadeN_ttf)

	var font = truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	text.Draw(dst, "Welcome to GoFract "+appVersion+", press 'h' for help & controls", font, 15+2, 28+2, color.Black)
	text.Draw(dst, "Welcome to GoFract "+appVersion+", press 'h' for help & controls", font, 15, 28, color.White)

	welcomeCounter--
}

//
// Debugging info
//
func drawDebugOverlay(dst *ebiten.Image) {
	var debugText = []string{
		fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()),
		fmt.Sprintf("magFactor: %0.8f maxIter: %.2f", f.MagFactor, f.MaxIter),
		fmt.Sprintf("center: %0.4f, %0.4f", f.Center.R, f.Center.I),
		fmt.Sprintf("juliaseed: %0.4f, %0.4f", f.JuliaSeed.R, f.JuliaSeed.I),
		fmt.Sprintf("renderTime: %vms", lastRenderTime),
		fmt.Sprintf("blendMode: %v", gradient.Mode),
	}
	var fontSize = int(float64(f.ImgWidth) / 60.0)
	var ttf, _ = truetype.Parse(exampleFonts.MPlus1pRegular_ttf)

	var font = truetype.NewFace(ttf, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	for ix, textString := range debugText {
		text.Draw(dst, textString, font, 5+2, ((ix+1)*(fontSize+6))+2, color.Black)
		text.Draw(dst, textString, font, 5, (ix+1)*(fontSize+6), color.White)
	}
}

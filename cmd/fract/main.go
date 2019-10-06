package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/benc-uk/gofract/pkg/fractals"
	"github.com/benc-uk/gofract/pkg/colors"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const appVersion = "0.0.3"

var (
	mainImage 	*image.RGBA
	gradient    colors.GradientTable
	showDebug 	bool
	showHelp  	bool
	lastRenderTime float64
	
	f fractals.Fractal // Global fractal super object
)

func update(screen *ebiten.Image) error {
	// Zoom in/out with mousewheel
	_, mouseDy := ebiten.Wheel()
	if mouseDy > 0 {
		f.MagFactor *= 0.8
		lastRenderTime = f.Render(mainImage, gradient)
	}
	if mouseDy < 0 {
		f.MagFactor *= 1.2
		lastRenderTime = f.Render(mainImage, gradient)
	}

	// On click recenter
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := ((float64(x) / float64(mainImage.Bounds().Max.X)) * f.W) - (f.W / 2.0)
		fy := ((float64(y) / float64(mainImage.Bounds().Max.Y)) * f.H) - (f.H / 2.0)
		f.Center.R += (fx * f.MagFactor)
		f.Center.I += (fy * f.MagFactor)
		lastRenderTime = f.Render(mainImage, gradient)
	}

	// 'S' key -> Save to PNG
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		filename := "Fractal_" + time.Now().Format("2006-01-02_150405") + ".png"
		fmt.Printf("Saving image to %v...\n", filename)
		f, _ := os.Create(filename)
		png.Encode(f, mainImage)
	}

	// 'D' key -> Enable debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		showDebug = !showDebug
		showHelp = false
	}

	// 'B' key -> cycle colour blend mode
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		gradient.Mode++
		gradient.Mode = gradient.Mode % colors.MaxColorModes
		lastRenderTime = f.Render(mainImage, gradient)
	}

	// 'M key -> create a random color palette
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		gradient.Randomise()
		lastRenderTime = f.Render(mainImage, gradient)
	}

	// 'R' key -> reload
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		oldWidth := f.ImgWidth
		initFractal()
		if(f.ImgWidth != oldWidth) {
			fmt.Println("### WHOA! Changing window size not supported, changes ignored until a restart, sorry!")
			f.ImgWidth = oldWidth
		}

		lastRenderTime = f.Render(mainImage, gradient)
	}


	// 'H' key -> show help
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		showHelp = !showHelp
		showDebug = false
	}

	// 'O' key -> show help
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {

	}

	// 'ESC' or 'Q' key -> quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	// IsRunningSlowly helps smooth things when it's running slow
	if !ebiten.IsRunningSlowly() {
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			f.JuliaSeed.R += 0.005 * (f.MagFactor/4.0)
			lastRenderTime = f.Render(mainImage, gradient)
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			f.JuliaSeed.R -= 0.005 * (f.MagFactor/4.0)
			lastRenderTime = f.Render(mainImage, gradient)
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			f.JuliaSeed.I -= 0.005 * (f.MagFactor/4.0)
			lastRenderTime = f.Render(mainImage, gradient)
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			f.JuliaSeed.I += 0.005 * (f.MagFactor/4.0)
			lastRenderTime = f.Render(mainImage, gradient)
		}

		if ebiten.IsKeyPressed(ebiten.KeyX) {
			f.MaxIter += 10
			lastRenderTime = f.Render(mainImage, gradient)
		}
		if ebiten.IsKeyPressed(ebiten.KeyZ) && f.MaxIter > 10 {
			f.MaxIter -= 10
			lastRenderTime = f.Render(mainImage, gradient)
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Render the offscreen image
	screen.ReplacePixels(mainImage.Pix)

	// Conditional text overlays
	if showDebug {
		drawDebugOverlay(screen)
	}
	if showHelp {
		drawHelpOverlay(screen)
	}
	if welcomeCounter > 0 {
		drawWelcomeOverlay(screen)
	}

	return nil
}

//
// Main entry point
//
func main() {
	fmt.Println("### Starting GoFract v" + appVersion + "...")

	initFractal()

	imgHeight := int(float64(f.ImgWidth) * float64(f.H / f.W))
	mainImage = image.NewRGBA(image.Rect(0, 0, f.ImgWidth, imgHeight))
	lastRenderTime = f.Render(mainImage, gradient)

	// Set icon (held in icon.go)
	icon, err := png.Decode(getIcon())
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	}
	ebiten.SetFullscreen(f.FullScreen)

	// It starts here
	fmt.Printf("### Window size: [%v,%v]\n", f.ImgWidth, imgHeight)
	if err := ebiten.Run(update, f.ImgWidth, imgHeight, 1, "GoFract v"+appVersion); err != nil {
		fmt.Println(err)
	}
}

func initFractal() {
	// Default fractal
	f = fractals.Fractal{
		FractType:  	"mandelbrot",
		Center:     	fractals.ComplexPair{-0.6, 0.0},
		MagFactor:  	1.0,
		MaxIter:    	90,
		W:         	 	3.0,
		H:         	 	2.0,
		ImgWidth: 	  1000,
		JuliaSeed:   	fractals.ComplexPair{0.355, 0.355},
		InnerColor: 	"#000000",
		FullScreen: 	false,
		ColorRepeats: 2.0,
	}

	// Handle loading YAML config file
	configFile := "fractal.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
		fmt.Println("### Trying to load:", configFile)
	}
	_, err := os.Stat(configFile)
	if err == nil {
		fmt.Println("### Loading config file:", configFile)
		fractals.LoadFractal(&f, configFile)
	} else {
		fmt.Println("### No config file, starting with defaults")
	}

	// Color gradient table
	if len(f.Colors) < 2 {
		gradient = colors.GradientTable{}
		gradient.AddToTable("#000762", 0.0)
		gradient.AddToTable("#0B48C3", 0.2)
		gradient.AddToTable("#ffffff", 0.4)
		gradient.AddToTable("#E3A000", 0.5)
		gradient.AddToTable("#000762", 0.9)
	} else {
		gradient = colors.GradientTable{}
		for _, col := range f.Colors {
			gradient.AddToTable(col.Color, col.Pos)
		}
	}
}
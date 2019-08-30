package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var (
	mainImage 	*image.RGBA
	gradient  	gradientTable
	debug     	= false

	renderTime 	float64

	imgWidth   	int
	imgHeight  	int
	imgWidthF  	float64
	imgHeightF 	float64

	// Global fractal super object
	f Fractal
)

type Fractal struct {
	fractType 					string
	magFactor 					float64
	centerX, centerY  	float64
	maxIter   					int
	c 									complex128

	w, h, ratioHW, ratioWH      float64
}

func renderFractal() {
	var wg sync.WaitGroup
	wg.Add(imgHeight)

	start := time.Now()
	// possible scaling to mag here
	//maxIter = int(math.Min(1000, 40.0 + (40.0 / magFactor)))

	for y := 0; y < imgHeight; y++ {
		// Use an anonymous goroutine to speed things up A LOT
		go func(y int) {
			for x := 0; x < imgWidth; x++ {
				// This gibberish converts from image space (x, y) to complex plane (r, i)
				// Takes into account aspect ratio, magnification and centering
				rOffset := f.centerX - (f.w / 2.0) * f.magFactor
				iOffset := f.centerY - (f.h / 2.0) * f.magFactor
				r := rOffset + ((float64(x) / imgWidthF) * f.w) * f.magFactor
				i := iOffset + ((float64(y) / imgHeightF) * f.h) * f.magFactor

				var iter int
				switch f.fractType {
				case "mandelbrot":
					iter = mandlebrot(complex(r, i), f)
				case "julia":
					iter = julia(complex(r, i), f)
				}

				// Output pixel colour, default is black
				pixelR := uint8(0)
				pixelG := uint8(0)
				pixelB := uint8(0)

				// Colour the pixel if it escaped, based on iteration count
				if iter < f.maxIter {
					scaledIter := float64(iter) / float64(f.maxIter)
					pixelR, pixelG, pixelB = gradient.getInterpolatedColorFor(scaledIter).RGB255()
				}

				// Store the pixel in the image buffer
				p := 4 * (x + y*imgWidth)
				mainImage.Pix[p] = pixelR
				mainImage.Pix[p+1] = pixelG
				mainImage.Pix[p+2] = pixelB
				mainImage.Pix[p+3] = 0xff
			}
			defer wg.Done()
		}(y)
	}

	wg.Wait()
	renderTime = float64(float64(time.Since(start)) / float64(time.Millisecond))
}

func update(screen *ebiten.Image) error {

	// Zoom in/out with mousewheel
	_, dy := ebiten.Wheel()
	if dy > 0 {
		f.magFactor *= 0.8
		renderFractal()
	}
	if dy < 0 {
		f.magFactor *= 1.2
		renderFractal()
	}

	// On click recenter
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := ((float64(x) / imgWidthF) * f.w) - (f.w / 2.0)
		fy := ((float64(y) / imgHeightF) * f.h) - (f.h / 2.0)
		f.centerX += (fx * f.magFactor)
		f.centerY += (fy * f.magFactor)
		renderFractal()
	}

	// 'S' key -> Save to PNG
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		filename := "Fractal_" + time.Now().Format("2006-01-02_15:04:05") + ".png"
		fmt.Printf("Saving image to %v...\n", filename)
		f, _ := os.Create(filename)
		png.Encode(f, mainImage)
	}

	// 'D' key -> Enable debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		debug = !debug
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Render the offscreen image
	screen.ReplacePixels(mainImage.Pix)
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("magFactor: %0.4f", f.magFactor), 2, 16)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Location: %0.4f, %0.4f", f.centerX, f.centerY), 2, 33)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time: %v  %v", renderTime, f.maxIter), 2, 49)
	}
	return nil
}

//
// Main entry point
//
func main() {
	fmt.Println("### Starting GoFract...")

	f = Fractal{
		fractType: "mandlebrot",
		centerX: 0.0,
		centerY: 0.0,
		magFactor: 1.0,
		maxIter: 100,
		w: 3.0,
		h: 2.0,
	}
	f.ratioHW = f.h / f.w
	f.ratioWH = f.w / f.h

	cr, ci := 0.0, 0.0
	flag.IntVar(&imgWidth, "width", 1000, "Image width")
	flag.IntVar(&f.maxIter, "maxiter", 100, "Max iterations")
	flag.StringVar(&f.fractType, "type", "mandelbrot", "Type of fractal")
	flag.Float64Var(&cr, "cr", 0.355, "Real value of c")
	flag.Float64Var(&ci, "ci", 0.355, "Imaginary value of c")
	flag.Parse()

	imgHeight = int(float64(imgWidth) * f.ratioHW)
	imgWidthF = float64(imgWidth)
	imgHeightF = float64(imgHeight)

	if f.fractType == "mandelbrot" {
		f.centerX = -0.6
	}

	f.c = complex(cr, ci)

	// Colour gradient table
	gradient = gradientTable{
		{parseHex("#010230"), 0.0},
		{parseHex("#090cb5"), 0.1},
		{parseHex("#7627ab"), 0.2},
		{parseHex("#f56320"), 0.3},
		{parseHex("#ffff00"), 1.0},
	}

	mainImage = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	renderFractal()

	icon, err := png.Decode(getIcon())
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	}

	fmt.Printf("### Window size: [%v,%v]\n", imgWidth, imgHeight)
	if err := ebiten.Run(update, imgWidth, imgHeight, 1, "GoFract"); err != nil {
		fmt.Println(err)
	}
}

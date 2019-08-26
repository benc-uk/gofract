package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	// "github.com/hajimehoshi/ebiten/inpututil"
)

var (
	mainImage  *image.RGBA
	gradient   gradientTable

	magFactor  = 1.0
	centerX    = -0.6
	centerY    = 0.0
	maxIter    = 80

	h          = 2.0
	w          = 3.0
	ratioHW    = h / w
	ratioWH    = w / h
)

var (
	imgWidth int
	imgHeight int
	imgWidthF float64
	imgHeightF float64
)

func renderFractal() {
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {

			// This gibberish converts from image space (x, y) to complex plane (r, i)
			// Takes into account aspect ratio, magnification and centering
			rOffset := centerX - ratioWH*magFactor
			iOffset := centerY - (h/2.0)*magFactor
			r := rOffset + ((float64(x)/imgWidthF)*w)*magFactor
			i := iOffset + ((float64(y)/imgHeightF)*h)*magFactor
			iter := mandlebrot(float64(r), float64(i))

			// Output pixel colour, default is black
			pixelR := uint8(0)
			pixelG := uint8(0)
			pixelB := uint8(0)

			// Colour the pixel if it escaped, based on iteration count
			if iter != maxIter {
				scaledIter := float64(iter) / float64(maxIter)
				pixelR, pixelG, pixelB = gradient.getInterpolatedColorFor(scaledIter).RGB255()
			}

			// Store the pixel in the image buffer
			p := 4 * (x + y*imgWidth)
			mainImage.Pix[p] = pixelR
			mainImage.Pix[p+1] = pixelG
			mainImage.Pix[p+2] = pixelB
			mainImage.Pix[p+3] = 0xff
		}
	}
}

func update(screen *ebiten.Image) error {

	// Zoom in/out with mousewheel
	_, dy := ebiten.Wheel()
	if dy > 0 {
		magFactor *= 0.8
		renderFractal()
	}
	if dy < 0 {
		magFactor *= 1.2
		renderFractal()
	}

	// On click recenter
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := ((float64(x) / imgWidthF) * w) - ratioWH
		fy := ((float64(y) / imgHeightF) * h) - (h / 2.0)
		centerX += (fx * magFactor)
		centerY += (fy * magFactor)
		renderFractal()
	}

	// 'S' key -> Save to PNG
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		filename := "Fractal_" + time.Now().Format("2006-01-02_15:04:05") + ".png"
		fmt.Printf("Saving image to %v...\n", filename)
		f, _ := os.Create(filename)
		png.Encode(f, mainImage)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Render the offscreen image
	screen.ReplacePixels(mainImage.Pix)
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	return nil
}

//
// Main entry point
//
func main() {
	flag.IntVar(&imgWidth, "width", 1000, "Image width")
	flag.IntVar(&maxIter, "maxiter", 80, "Max iterations")
	flag.Parse()
	imgHeight  = int(float64(imgWidth) * ratioHW)
	imgWidthF  = float64(imgWidth)
	imgHeightF = float64(imgHeight)

	// Colour gradient table
	gradient = gradientTable{
		{parseHex("#090cb5"), 0.0},
		{parseHex("#7627ab"), 0.1},
		{parseHex("#f56320"), 0.3},
		{parseHex("#ffff00"), 1.0},
	}

	mainImage = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	renderFractal()

	icon, err := png.Decode(getIcon())
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	}

	fmt.Println("### Starting GoFract...")
	fmt.Printf("### Window size: [%v,%v]\n", imgWidth, imgHeight)
	if err := ebiten.Run(update, imgWidth, imgHeight, 1, "GoFract"); err != nil {
		fmt.Println(err)
	}
}

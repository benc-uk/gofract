package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
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
	imgWidth   = 1000
	imgHeight  = int(float64(imgWidth) * ratioHW)
	imgWidthF  = float64(imgWidth)
	imgHeightF = float64(imgHeight)
)

func renderFractal() {
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			rOffset := centerX - ratioWH*magFactor
			iOffset := centerY - (h/2.0)*magFactor
			r := rOffset + ((float64(x)/imgWidthF)*w)*magFactor
			i := iOffset + ((float64(y)/imgHeightF)*h)*magFactor
			iter := mandlebrot(float64(r), float64(i))

			rc := uint8(0)
			gc := uint8(0)
			bc := uint8(0)
			if iter != maxIter {
				scaledIter := float64(iter) / float64(maxIter)
				rc, gc, bc = gradient.getInterpolatedColorFor(scaledIter).RGB255()
			}

			p := 4 * (x + y*imgWidth)
			mainImage.Pix[p] = rc
			mainImage.Pix[p+1] = gc
			mainImage.Pix[p+2] = bc
			mainImage.Pix[p+3] = 0xff
		}
	}
}

func update(screen *ebiten.Image) error {
	_, dy := ebiten.Wheel()
	if dy > 0 {
		magFactor *= 0.8
		renderFractal()
	}
	if dy < 0 {
		magFactor *= 1.2
		renderFractal()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := ((float64(x) / imgWidthF) * w) - ratioWH
		fy := ((float64(y) / imgHeightF) * h) - (h / 2.0)
		centerX += (fx * magFactor)
		centerY += (fy * magFactor)
		renderFractal()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		fmt.Println("Saving image to PNG...")
		f, _ := os.Create("fractal_" + time.Now().String() + ".png")
		png.Encode(f, mainImage)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.ReplacePixels(mainImage.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	return nil
}

func init() {
	gradient = gradientTable{
		{parseHex("#090cb5"), 0.0},
		{parseHex("#7627ab"), 0.1},
		{parseHex("#f56320"), 0.3},
		{parseHex("#ffff00"), 1.0},
	}

	mainImage = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	renderFractal()
}

func main() {
	fmt.Println("### Starting... ", imgWidth, imgHeight)
	if err := ebiten.Run(update, imgWidth, imgHeight, 1, "GoFract"); err != nil {
		fmt.Println(err)
	}
}

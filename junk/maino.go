package maino

import (
	// "fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	imgWidth := 1600
	imgHeight := 1200
	magFactor := 500.0
	panX := 2.4
	panY := 1.2

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imgWidth, imgHeight}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	pink := color.RGBA{255, 25, 200, 0xff}

	for x := 0; x < imgWidth; x++ {
		for y := 0; y < imgHeight; y++ {
			belongsToSet := checkInSet(float64(float64(x)/magFactor - panX), float64(float64(y)/magFactor - panY));
			if(belongsToSet > 0.0) {
				img.Set(x, y, pink)
			} else {
				img.Set(x, y, color.Black)
			}          
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}



func checkInSet(r, i float64) float64 {
	realComponentOfResult := r;
	imaginaryComponentOfResult := i;
	
	for iter := 0; iter < 50; iter++ {
		// Calculate the real and imaginary components of the result
		// separately
		tempRealComponent := (realComponentOfResult * realComponentOfResult - imaginaryComponentOfResult * imaginaryComponentOfResult + r)

		tempImaginaryComponent := (2.0 * realComponentOfResult * imaginaryComponentOfResult + i)

		realComponentOfResult = tempRealComponent;
		imaginaryComponentOfResult = tempImaginaryComponent;
	}
	
	if (realComponentOfResult * imaginaryComponentOfResult < 5) {
		return (realComponentOfResult * imaginaryComponentOfResult); // In the Mandelbrot set
	}

	return 0.0; // Not in the set
}

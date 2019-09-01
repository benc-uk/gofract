package fractals

import (
	"fmt"
	"image"
	"io/ioutil"
	"math"
	//"math/cmplx"
	"os"
	"sync"
	"time"

	"github.com/benc-uk/gofract/pkg/colors"

	yaml "gopkg.in/yaml.v2"
)

var (
	escape = 256.0
	escape2 = escape * escape
	log2   = math.Log(2.0)
)

// Render the fractal into the given image using the given palette 
func (f Fractal) Render(img *image.RGBA, palette colors.GradientTable) float64 {
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y

	var seed complex128
	if f.FractType == "julia" {
		seed = complex(f.JuliaSeed.R, f.JuliaSeed.I)
	}

	innerColor := colors.ParseHex(f.InnerColor)
	innerR := uint8(innerColor.R * 255)
	innerG := uint8(innerColor.G * 255)
	innerB := uint8(innerColor.B * 255)

	var wg sync.WaitGroup
	wg.Add(imgHeight)

	start := time.Now()
	for y := imgHeight-1; y >= 0 ; y-- {
		// Use an anonymous goroutine to speed things up A LOT
		go func(y int) {
			for x := 0; x < imgWidth; x++ {
				// This gibberish converts from image space (x, y) to complex plane (r, i)
				// Takes into account aspect ratio, magnification and centering
				rOffset := f.Center.R - (f.W/2.0)*f.MagFactor
				iOffset := f.Center.I - (f.H/2.0)*f.MagFactor
				r := rOffset + ((float64(x)/float64(imgWidth))*f.W)*f.MagFactor
				i := iOffset + ((float64(y)/float64(imgHeight))*f.H)*f.MagFactor

				var iter float64
				switch f.FractType {
					case "mandelbrot":
						iter = mandlebrot(complex(r, i), f)
					case "julia":
						iter = julia(complex(r, i), f, seed)
					default:
						iter = mandlebrot(complex(r, i), f)
				}

				// Default to inner colour if inside the set
				pixelR, pixelG, pixelB := innerR, innerG, innerB

				// Color the pixel if it escaped, based on iteration count
				if iter < f.MaxIter {
					// This maths lets us have repeating colors
					repeatSize := f.MaxIter / f.ColorRepeats
					scaledIter := math.Mod(iter, repeatSize) / repeatSize
					pixelR, pixelG, pixelB = palette.GetInterpolatedColorFor(scaledIter).RGB255()
				}

				// Store the pixel in the image buffer
				p := 4 * (x + y*f.ImgWidth)
				img.Pix[p] = pixelR
				img.Pix[p+1] = pixelG
				img.Pix[p+2] = pixelB
				img.Pix[p+3] = 0xff
			}
			defer wg.Done()
		}(y)
	}

	wg.Wait()
	return float64(float64(time.Since(start)) / float64(time.Millisecond))
}

func mandlebrot(a complex128, f Fractal) float64 {
	var z complex128 // zero
	var iter float64
	var mag float64
	for iter <= f.MaxIter {
		z = z*z + a
		mag = real(z)*real(z)+imag(z)*imag(z)
		if mag > escape2 {
			break
		}
		iter++
	}

	if iter >= f.MaxIter {
		return float64(f.MaxIter)
	}

	// I have NO IDEA if this is correct but it looks good
	smoothIter := iter - math.Log(math.Log(mag/math.Log(escape)))/log2
	return smoothIter
}

func julia(a complex128, f Fractal, seed complex128) float64 {
	z := a
	var iter float64
	var mag float64
	for iter <= f.MaxIter {
		z = z*z + seed
		mag = real(z)*real(z)+imag(z)*imag(z)
		if mag > escape2 {
			break
		}
		iter++
	}

	if iter >= f.MaxIter {
		return float64(f.MaxIter)
	}

	// I have NO IDEA if this is correct but it looks good
	smoothIter := iter - math.Log(math.Log(mag/math.Log(escape)))/log2
	return smoothIter
}

// LoadFractal is YAML parser and loader
func LoadFractal(f *Fractal, filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("### Error loading YAML, %v", err)
		os.Exit(1)
	}

	err = yaml.UnmarshalStrict(yamlFile, f)
	if err != nil {
		fmt.Printf("### Error unmarshalling YAML, %v", err)
		os.Exit(2)
	}
}

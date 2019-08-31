package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	yaml "gopkg.in/yaml.v2"
)

var (
	mainImage *image.RGBA
	gradient  gradientTable
	debug     = false

	renderTime float64

	// Global fractal super object
	f     Fractal
)

type ComplexPair struct {
	R float64 `yaml:"r"`
	I float64 `yaml:"i"`
}

type ColorDef struct {
	Pos   float64 `yaml:"pos"`
	Color string  `yaml:"color"`
}

type Fractal struct {
	FractType  string      `yaml:"type"`
	MagFactor  float64     `yaml:"zoom"`
	MaxIter    int         `yaml:"maxIter"`
	W          float64     `yaml:"width"`
	H          float64     `yaml:"height"`
	ImgWidth   int         `yaml:"imageWidth"`
	Center     ComplexPair `yaml:"center"`
	JuliaC     ComplexPair `yaml:"juliaC"`
	Colors     []ColorDef  `yaml:"colors"`
	InnerColor string      `yaml:"innerColor"`
	FullScreen bool        `yaml:"fullScreen"`

	imgHeight int
	ratioHW   float64
	ratioWH   float64
}

func renderFractal(fast bool) {
	imgWidthF := float64(f.ImgWidth)
	imgHeightF := float64(f.imgHeight)

	var c complex128
	if f.FractType == "julia" {
		c = complex(f.JuliaC.R, f.JuliaC.I)
	}

	innerColor := parseHex(f.InnerColor)
	innerR := uint8(innerColor.R * 255)
	innerG := uint8(innerColor.G * 255)
	innerB := uint8(innerColor.B * 255)

	var wg sync.WaitGroup
	wg.Add(f.imgHeight)

	start := time.Now()
	// possible scaling to mag here
	//maxIter = int(math.Min(1000, 40.0 + (40.0 / magFactor)))

	for y := 0; y < f.imgHeight; y++ {
		// Use an anonymous goroutine to speed things up A LOT
		go func(y int) {
			for x := 0; x < f.ImgWidth; x++ {
				// This gibberish converts from image space (x, y) to complex plane (r, i)
				// Takes into account aspect ratio, magnification and centering
				rOffset := f.Center.R - (f.W/2.0)*f.MagFactor
				iOffset := f.Center.I - (f.H/2.0)*f.MagFactor
				r := rOffset + ((float64(x)/imgWidthF)*f.W)*f.MagFactor
				i := iOffset + ((float64(y)/imgHeightF)*f.H)*f.MagFactor

				var iter int
				switch f.FractType {
					case "mandelbrot":
						iter = mandlebrot(complex(r, i), f)
					case "julia":
						iter = julia(complex(r, i), f, c)
					default:
						iter = mandlebrot(complex(r, i), f)
				}

				// Default to inner colour if inside the set
				pixelR, pixelG, pixelB := innerR, innerG, innerB

				// Color the pixel if it escaped, based on iteration count
				if iter < f.MaxIter {
					scaledIter := float64(iter) / float64(f.MaxIter)
					pixelR, pixelG, pixelB = gradient.getInterpolatedColorFor(scaledIter).RGB255()
				}

				// Store the pixel in the image buffer
				p := 4 * (x + y*f.ImgWidth)
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
	_, mouseDy := ebiten.Wheel()
	if mouseDy > 0 {
		f.MagFactor *= 0.8
		renderFractal(true)
	}
	if mouseDy < 0 {
		f.MagFactor *= 1.2
		renderFractal(true)
	}

	// On click recenter
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fx := ((float64(x) / float64(f.ImgWidth)) * f.W) - (f.W / 2.0)
		fy := ((float64(y) / float64(f.imgHeight)) * f.H) - (f.H / 2.0)
		f.Center.R += (fx * f.MagFactor)
		f.Center.I += (fy * f.MagFactor)
		renderFractal(true)
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
		debug = !debug
	}

	// 'B' key -> cycle colour blend mode
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		gradient.mode++
		gradient.mode = gradient.mode % maxColorModes
		renderFractal(true)
	}

	// 'R' key -> cycle colour blend mode
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		gradient.randomise()
		renderFractal(true)
	}

	// 'ESC' or 'Q' key -> quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		f.JuliaC.R += 0.005 * f.MagFactor
		renderFractal(true)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		f.JuliaC.R -= 0.005 * f.MagFactor
		renderFractal(true)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		f.JuliaC.I -= 0.005 * f.MagFactor
		renderFractal(true)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		f.JuliaC.I += 0.005 * f.MagFactor
		renderFractal(true)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Render the offscreen image
	screen.ReplacePixels(mainImage.Pix)

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\n", ebiten.CurrentTPS()))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("magFactor: %0.4f maxIter: %d", f.MagFactor, f.MaxIter), 2, 16)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("center: %0.4f, %0.4f", f.Center.R, f.Center.I), 2, 33)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("juliaC: %0.4f, %0.4f", f.JuliaC.R, f.JuliaC.I), 2, 49)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("renderTime: %vms", renderTime), 2, 65)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("blendMode: %v", gradient.mode), 2, 81)
	}

	return nil
}

//
// Main entry point
//
func main() {
	fmt.Println("### Starting GoFract...")

	f = Fractal{
		FractType:  "mandelbrot",
		Center:     ComplexPair{-0.6, 0.0},
		MagFactor:  1.0,
		MaxIter:    80,
		W:          3.0,
		H:          2.0,
		ImgWidth:   1000,
		JuliaC:     ComplexPair{0.355, 0.355},
		InnerColor: "#000000",
		FullScreen: false,
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
		f.loadFractal(configFile)
	} else {
		fmt.Println("### No config file, starting with defaults")
	}

	f.ratioHW = f.H / f.W
	f.ratioWH = f.W / f.H
	f.imgHeight = int(float64(f.ImgWidth) * f.ratioHW)

 	fractalYamlDump, _ := yaml.Marshal(f)
	fmt.Printf("\n%v\n", string(fractalYamlDump))

	// Color gradient table
	if len(f.Colors) < 2 {
		gradient = gradientTable{}
		gradient.addToTable("#010230", 0.0)
		gradient.addToTable("#090cb5", 0.1)
		gradient.addToTable("#7627ab", 0.2)
		gradient.addToTable("#f56320", 0.3)
		gradient.addToTable("#ffff00", 1.0)
	} else {
		gradient = gradientTable{}
		for _, col := range f.Colors {
			gradient.addToTable(col.Color, col.Pos)
		}
	}

	mainImage = image.NewRGBA(image.Rect(0, 0, f.ImgWidth, f.imgHeight))
	renderFractal(false)

	// Set icon (held in icon.go)
	icon, err := png.Decode(getIcon())
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	}
	ebiten.SetFullscreen(f.FullScreen)

	// It starts here
	fmt.Printf("### Window size: [%v,%v]\n", f.ImgWidth, f.imgHeight)
	if err := ebiten.Run(update, f.ImgWidth, f.imgHeight, 1, "GoFract v0.0.3"); err != nil {
		fmt.Println(err)
	}
}

//
// YAML parser
//
func (f *Fractal) loadFractal(filename string) *Fractal {
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
	return f
}

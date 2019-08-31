# GoFract
Mandlebrot and Julia fractals rendered in real-time using Go. 

Uses the [Enbiten](https://ebiten.org/) "A dead simple 2D game library" and also [go-colorful](https://github.com/lucasb-eyer/go-colorful) a library for manipulating colours in Go.

It should build/run under both Linux/WSL and Windows

Features:
- Mandlebrot and Julia sets
- Zoom in/out with mouse wheel
- Click anywhere to recenter the view
- Press 's' to save current view to a PNG file
- Press 'd' to display debugging information
- Press 'r' to randomize the colour pallet
- Press 'b' to change the colour blend mode (RGB, HCL, HSV)
- Use cursor keys to explore when in Julia set mode, by changing the real/imaginary parts of C

## Configuration
Configuration is done via YAML, this is loaded from `fractal.yaml` by default, or the filename can be passed as an argument when starting the app.

Below is an example config file, not every setting needs to be provided and if the file is not provided defaults are used.

```yaml
type: julia         # Either 'mandelbrot' or 'julia', default: mandelbrot
maxIter: 200        # Max iterations, default: 80
imageWidth: 800     # Width of the windows and image, default: 1000
fullScreen: false   # Run fullscreen, default: false

width: 3.0          # Width in the complex plane (real part), default: 3.0
height: 2.0         # Height in the complex plane (imaginary part), default: 2.0
                    # NOTE. The ratio of width:height combined with imageWidth defines the imageHeight
                    #       imageHeight = imageWidth * (height / width)
zoom: 1.5           # Starting zoom factor, default: 1.0

center:             # Starting location in complex plane, default: [0.0,-0.6]
  r: 0.0
  i: -0.6

juliaC:             # Used when type=julia, complex C value used, default: [0.355, 0.355]
  r: -0.54
  i: 0.54

# Array of colors (in hex format), and positions, pos: 0.0 ~ 1.0 
# Minimum of two colors, colors are blended to make a smooth gradient between 0.0 and 1.0
colors:              
  - pos: 0.0
    color: "#130b5c"
  - pos: 0.2
    color: "#3ec71c"
  - pos: 0.7
    color: "#db4918"         
  - pos: 1.0
    color: "#cf0c84"      

innerColor: "#570336"   # Color used to draw inside the fractal set, default = #000000
```

## Screen shots
#### [Gallery Here](https://code.benco.io/gofract/img/)

---

# Building Yourself
Tested with Go 1.12

## Linux 

Install prereq libs
```
sudo apt install libgl1-mesa-dev xorg-dev
```

Run directly
```
cd cmd/fract
go run .
```

Or build exe
```
cd cmd/fract
go build
```

## Windows
- Have Go installed / on your path
- Have Git installed / on your path
- Install mingw-w64 http://win-builds.org/doku.php/download_and_installation_from_windows and put bin directory on your path
- Build or run as with Linux above

---

# Appendix - X11 on WSL

Getting X11 windowing system and GUI working with WSL is "interesting"...

In WSL install xfce4, there might be other ways to set up X11 and associated libraries, but this is a lightweight windowing manager which meets most needs
```
sudo apt install xfce4 xfce4-terminal
```

In Windows download and install [VcXsrv](https://sourceforge.net/projects/vcxsrv/)

When starting VcXsrv **YOU MUST** set these extra settings:
- Un-check: "Native opengl"
- Check: "Disable access control"


Export display variable, trick here is to use the IP assigned to your "main" (LAN/Wifi) adapter on the Windows side, e.g.
```
export DISPLAY=192.168.0.24:0
```

Make sure `LIBGL_ALWAYS_INDIRECT` is NOT set

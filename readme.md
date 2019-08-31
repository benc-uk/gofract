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
- Press 'r' to randomise the colour pallet
- Press 'b' to change the colour blend mode (RGB, HCL, HSV)
- Use cursor keys to explore when in Julia set mode, by changing the real/imaginary parts of C

Command line options:
- `-width` Windows width in pixels (default 1000)
- `-maxiter` Maximum fractal iteration per pixel (default 80)
- `-type` Fractal type "mandelbrot" or "julia" (default mandelbrot)
- `-cr` When in Julia mode, set the real part of C (default 0.355)
- `-ci` When in Julia mode, set the imaginary part of C (default 0.355)
- `-colors` Specify a gradient colour palette, in the form of a comma separated list of pairs `hexcolor1=pos1,hexcolor2=pos2`. Where `hexcolor` is a colour in HEX form, e.g ff34bb (without hash), and `pos` is a position, between 0.0 and 1.0

## Screen shots
### [Gallery Here](https://code.benco.io/gofract/img/)

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

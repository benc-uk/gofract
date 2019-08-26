# GoFract
Fractals rendered in Go. Uses the [Enbiten](https://ebiten.org/) "A dead simple 2D game library"

It should run under Linux/WSL and Windows

Features:
- Zoom in/out with mouse wheel
- Click to recenter
- Press 's' to save current view to a PNG file

Command Line:
- `-width` Windows width in pixels (default 1000)
- `-maxiter` Maximum fractal iteration per pixel (default 80)

## Screen shots
![](https://user-images.githubusercontent.com/14982936/63654094-b09db900-c76d-11e9-90a3-e4540944f17c.png)

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

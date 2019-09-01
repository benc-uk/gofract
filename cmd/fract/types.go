package main

// ComplexPair is used to represent complex number in YAML
type ComplexPair struct {
	R float64 `yaml:"r"`
	I float64 `yaml:"i"`
}

// ColorDef is used to represent a gradient colour postion in YAML
type ColorDef struct {
	Pos   float64 `yaml:"pos"`
	Color string  `yaml:"color"`
}

// Fractal is our main object
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
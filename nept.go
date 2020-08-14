package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
	_ "image/gif"
	"math/rand"
	"os"
  "time"
)

// Calculation constants
const (
	rgbMax     float32 = 255.0
	percentMax float32 = 100.0
	intMax     float32 = 65535
)

// Global Flag definitions
var in = flag.String("i", "", "input file to adjust")
var out = flag.String("o", "", "output file to adjust")
var debug = flag.Bool("debug", false, "debugging on")

// Editing Flags
var bright = flag.Int("b", 0, "brighten")
var dark = flag.Int("d", 0, "darken")
var flat = flag.Int("f", 0, "flatten")
var iso = flag.Int("s", 0, "iso")

// Pixel
type Pixel struct {
	r, g, b, a uint32
}

func editPixel(x, y int, src image.Image, img *image.RGBA) {

  debugging("\nEditing Pixel %d:%d", x, y)

  // read values from original pixel and create new struct
  r, g, b, a := src.At(x, y).RGBA()
  pixel := Pixel{r: r, g: g, b: b, a: a}

  debugging("Original: %+v", pixel)

  if *bright > 0 {
    pixel = brighten(pixel, uint32(*bright))
  }

  if *dark > 0 {
    pixel = darken(pixel, uint32(*dark))
  }

  if *flat > 0 {
    pixel = flatten(pixel, uint32(*flat))
  }

  if *iso > 0 {
    pixel = isoify(pixel, uint32(*iso))
  }

  debugging("Modified: %+v", pixel)
  img.Set(x, y, constructRGBA(pixel))
}

func main() {

	flag.Parse()

  // Open File and read
  src := readImage(*in)

	// initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Fetch image dimensions
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	// initialize new image
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
      editPixel(x, y, src, img)
		}
	}

	// Encode as PNG.
	f, _ := os.Create(*out)
	png.Encode(f, img)

}

// Debugging Messages
func debugging(format string, args ...interface{}) {
	if *debug == true {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

// Reads various filetypes from file from argument
func readImage(in string) image.Image {
	infile, err := os.Open(in)
	if err != nil {
		panic("Please specify input file using -i")
	}
	defer infile.Close()


	src, _, err := image.Decode(infile)
	if err != nil {
		// replace this with real error handling
		panic(err)
	}

  return src
}

// Construct Pixel
func constructRGBA(p Pixel) color.RGBA {
	return color.RGBA{uint32ToRGB(p.r), uint32ToRGB(p.g), uint32ToRGB(p.b), 255}
}

// convert modifier % in points in rgb255 space
// ie. 20% -> 51 points
func percentToInt(val uint32) uint32 {
	debugging("percenttoInt: %d -> %d", val, uint32(intMax/percentMax*float32(val)))
	return uint32(intMax / percentMax * float32(val))
}

// convert input uint32 to rgb255
// ie. 60500 -> 235
func uint32ToRGB(val uint32) uint8 {
	debugging("uint32ToRGB: %d -> %d", val, uint8(rgbMax/percentMax*float32(val)/intMax*percentMax))
	return uint8(rgbMax / percentMax * float32(val) / intMax * percentMax)
}

// Adds values to a color and checks boundaries
// 67000 -> 65535
func addInt(c, v uint32) uint32 {
	o := c

	if c+v > uint32(intMax) {
		c = uint32(intMax)
	} else {
		c += v
	}

	debugging("addInt: %d + %d = %d", o, v, c)
	return c
}

// Subs values to a color and checks boundaries
// 4444232312 -> 0
func subInt(c, v uint32) uint32 {
	o := c

	if c < v {
		c = 0
	} else {
		c -= v
	}
	debugging("subInt: %d - %d = %d", o, v, c)

	return c
}

// Add a certain percentage of each rgb value
// [ (255,255,255), (145,77,83), ... ]
// ... to ...
// [ (235,235,235), (115,47,63), ... ]
func brighten(p Pixel, v uint32) Pixel {
	p.r = addInt(p.r, percentToInt(v))
	p.g = addInt(p.g, percentToInt(v))
	p.b = addInt(p.b, percentToInt(v))
	return p
}

// Reduce a certain percentage of each rgb value
// [ (255,255,255), (145,77,83), ... ]
// ... to ...
// [ (235,235,235), (115,47,63), ... ]
func darken(p Pixel, v uint32) Pixel {
	p.r = subInt(p.r, percentToInt(v))
	p.g = subInt(p.g, percentToInt(v))
	p.b = subInt(p.b, percentToInt(v))
	return p
}

// Raises lower imits of blacks what results in reduing depth and details.
// [ (48, 150, 30) ]
//    < ,  > , <
// [ (50, 150, 50) ]
func flatten(p Pixel, v uint32) Pixel {
	if p.r < percentToInt(v) {
		p.r = percentToInt(v)
	}

	if p.g < percentToInt(v) {
		p.g = percentToInt(v)
	}

	if p.b < percentToInt(v) {
		p.b = percentToInt(v)
	}

	return p
}

func isoify(p Pixel, v uint32) Pixel {
	r := uint32(rand.Int31n(int32(v)))
	p.r = addInt(p.r, percentToInt(r))
	p.g = addInt(p.g, percentToInt(r))
	p.b = addInt(p.b, percentToInt(r))
	return p
}

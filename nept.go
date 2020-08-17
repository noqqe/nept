package main

import (
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Calculation constants
const (
	rgbMax     float32 = 255.0
	percentMax float32 = 100.0
	intMax     float32 = 65535
)

// Pixel
type Pixel struct {
	r, g, b, a uint32
}

// Modification
type Modifications struct {
	x, y                    int  // Coords
	bright, dark, flat, iso int  // Mods on scale
	neg                     bool // Mods enable/disable
}

func editPixel(m Modifications, src image.Image, img *image.RGBA, wg *sync.WaitGroup) {

	// waiting group cancel after function
	defer wg.Done()

	// read values from original pixel and create new struct
	r, g, b, a := src.At(m.x, m.y).RGBA()
	pixel := Pixel{r: r, g: g, b: b, a: a}

	if m.neg == true {
		pixel = negative(pixel)
	}

	if m.bright > 0 {
		pixel = brighten(pixel, uint32(m.bright))
	}

	if m.dark > 0 {
		pixel = darken(pixel, uint32(m.dark))
	}

	if m.flat > 0 {
		pixel = flatten(pixel, uint32(m.flat))
	}

	if m.iso > 0 {
		pixel = isoify(pixel, uint32(m.iso))
	}

	img.Set(m.x, m.y, constructRGBA(pixel))
}

func main() {

	// Global Flag definitions
	var in string
	var out string
	var neg bool
	var bright, dark, flat, iso int
	var wg sync.WaitGroup

	// Option Parser
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "brightness",
				Value:       0,
				Usage:       "Brighten the Image",
				Aliases:     []string{"b"},
				Destination: &bright,
			},
			&cli.IntFlag{
				Name:        "darkness",
				Value:       0,
				Usage:       "Darken the image",
				Aliases:     []string{"d"},
				Destination: &dark,
			},
			&cli.IntFlag{
				Name:        "flatten",
				Value:       0,
				Usage:       "Flatten the image",
				Aliases:     []string{"f"},
				Destination: &flat,
			},
			&cli.IntFlag{
				Name:        "iso",
				Value:       0,
				Usage:       "Add iso to the image",
				Aliases:     []string{"s"},
				Destination: &iso,
			},
			&cli.BoolFlag{
				Name:        "negative",
				Value:       false,
				Usage:       "Convert negative to positive image",
				Aliases:     []string{"n"},
				Destination: &neg,
			},
			&cli.PathFlag{
				Name:        "in",
				Usage:       "Image to edit (input)",
				Aliases:     []string{"i"},
				Destination: &in,
				Required:    true,
				TakesFile:   true,
			},
			&cli.PathFlag{
				Name:        "out",
				Usage:       "Image to edit (output)",
				Aliases:     []string{"o"},
				Destination: &out,
				Required:    true,
				TakesFile:   true,
			},
		},
		Action: func(c *cli.Context) error {

			// Open File and read
			src := readImage(in)

			// initialize random seed
			rand.Seed(time.Now().UnixNano())

			// Fetch image dimensions
			bounds := src.Bounds()
			w, h := bounds.Max.X, bounds.Max.Y

			// Initialize Progress Bar
			bar := progressbar.NewOptions(w*h,
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetWidth(30),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}))

			// initialize new image
			img := image.NewRGBA(image.Rect(0, 0, w, h))

			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					wg.Add(1)
					bar.Add(1)

					go editPixel(
						Modifications{
							x:      x,
							y:      y,
							bright: bright,
							dark:   dark,
							flat:   flat,
							iso:    iso,
							neg:    neg},
						src, img, &wg)

				}
			}
			wg.Wait()

			// Encode as PNG.
			f, _ := os.Create(out)
			png.Encode(f, img)

			return nil
		},
	}

	app.Run(os.Args)
}

// Reads various filetypes from file from argument
func readImage(in string) image.Image {
	infile, err := os.Open(in)
	if err != nil {
		log.Fatal("No such file or directory: ", in)
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
	return uint32(intMax / percentMax * float32(val))
}

// convert input uint32 to rgb255
// ie. 60500 -> 235
func uint32ToRGB(val uint32) uint8 {
	return uint8(rgbMax / percentMax * float32(val) / intMax * percentMax)
}

// Adds values to a color and checks boundaries
// 67000 -> 65535
func addInt(c, v uint32) uint32 {
	if c+v > uint32(intMax) {
		c = uint32(intMax)
	} else {
		c += v
	}

	return c
}

// Subs values to a color and checks boundaries
// 4444232312 -> 0
func subInt(c, v uint32) uint32 {
	if c < v {
		c = 0
	} else {
		c -= v
	}

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

// Adds noise to an image resulting in having an high ISO like
// look. Done by adding random values in a certain range to the rgb value.
// [ (48, 150, 30), (49, 151, 31) ]
//    +10,+10,+10    +1, +1, +1
func isoify(p Pixel, v uint32) Pixel {
	r := uint32(rand.Int31n(int32(v)))
	p.r = addInt(p.r, percentToInt(r))
	p.g = addInt(p.g, percentToInt(r))
	p.b = addInt(p.b, percentToInt(r))
	return p
}

// Converts each value into its opposite
// Useful for converting negative film scans
// [ (2, 255, 250) ]
// [ (254, 0,		4) ]
func negative(p Pixel) Pixel {
	p.r = uint32(intMax) - p.r
	p.g = uint32(intMax) - p.g
	p.b = uint32(intMax) - p.b
	return p
}

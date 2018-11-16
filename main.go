package main

import (
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

// https://stackoverflow.com/questions/42516203/converting-rgba-image-to-grayscale-golang
func toGray16(img image.Image) *image.Gray16 {
	b := img.Bounds()
	ret := image.NewGray16(b)

	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			oldPixel := img.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			ret.Set(x, y, pixel)
		}
	}

	return ret
}

func main() {
	name := flag.String("f", "", "")
	flag.Parse()

	f, err := os.Open(*name)
	if err != nil {
		panic(err)
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		panic(err)
	}

	g := toGray16(img)

	w, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	png.Encode(w, g)
}

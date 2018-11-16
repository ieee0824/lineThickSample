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

func gaussian(img *image.Gray16) *image.Gray16 {
	b := img.Bounds()
	ret := image.NewGray16(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			var buf uint
			var counter uint
		skip:
			for fy := -2; fy <= 2; fy++ {
				for fx := -2; fx <= 2; fx++ {
					c := img.At(x+fy, y+fy).(color.Gray16)
					if y+fy < b.Min.Y || x+fx < b.Min.X {
						ret.Set(x, y, img.At(x, y))
						break skip
					}
					counter++
					buf += uint(c.Y)
				}
			}
			if counter == 0 {
				continue
			}
			ret.Set(x, y, color.Gray16{Y: uint16(buf / counter)})
		}
	}

	return ret
}

func fault(img *image.Gray16, t uint16) *image.Gray16 {
	b := img.Bounds()
	ret := image.NewGray16(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			ret.Set(x, y, color.Gray16{65535})
			if y+1 < b.Max.Y {
				p := img.At(x, y).(color.Gray16)
				q := img.At(x, y+1).(color.Gray16)
				n := int(p.Y) - int(q.Y)
				if n < 0 {
					n *= -1
				}
				if t < uint16(n) {
					ret.Set(x, y, color.Gray16{0})
				}
			}
			if x+1 < b.Max.X {
				p := img.At(x, y).(color.Gray16)
				q := img.At(x+1, y).(color.Gray16)
				n := int(p.Y) - int(q.Y)
				if n < 0 {
					n *= -1
				}
				if t < uint16(n) {
					ret.Set(x, y, color.Gray16{0})
				}
			}
		}
	}

	return ret
}

func edge(img *image.Gray16) *image.Gray16 {
	return fault(img, 512)
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

	g := edge(gaussian(toGray16(img)))

	w, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	png.Encode(w, g)
}

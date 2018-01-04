package jump

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func toImg(img [][]int) image.Image {
	limitX, limitY := len(img), len(img[0])
	pic := image.NewNRGBA(image.Rect(0, 0, limitX, limitY))

	for x := 0; x < limitX; x++ {
		for y := 0; y < limitY; y++ {
			if img[x][y] == 1 {
				pic.Set(x, y, color.NRGBA{
					R: uint8((0) & 255),
					G: uint8((0) << 1 & 255),
					B: uint8((0) << 2 & 255),
					A: 255,
				})
			} else {
				pic.Set(x, y, color.NRGBA{
					R: uint8((255) & 255),
					G: uint8((255) << 1 & 255),
					B: uint8((255) << 2 & 255),
					A: 255,
				})
			}
		}
	}
	return pic
}

func savePNG(fileName string, pic image.Image) {
	go func() {
		f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		png.Encode(f, pic)
		f.Close()
	}()
}

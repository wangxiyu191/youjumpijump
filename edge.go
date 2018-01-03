package jump

import (
	"image"
	"image/color"
	"math"
)

func EdgeDetection(pic image.Image, QE float64) image.Image {
	limit := pic.Bounds().Size()

	img := image.NewNRGBA(image.Rect(0, 0, limit.X, limit.Y))

	for y := 0; y < limit.Y; y++ {
		for x := 0; x < limit.X; x++ {

			var sum float64

			pixelO := pic.At(x, y)
			x2, y2, z2, _ := pixelO.RGBA()

			for _, offsetX := range []int{-1, 0, 1} {
				for _, offsetY := range []int{-1, 0, 1} {
					if offsetX == offsetY {
						break
					}
					pixelN := pic.At(x+offsetX, y+offsetY)
					x1, y1, z1, _ := pixelN.RGBA()

					xSqr := (x1 - x2) * (x1 - x2)
					ySqr := (y1 - y2) * (y1 - y2)
					zSqr := (z1 - z2) * (z1 - z2)
					mySqr := float64(xSqr + ySqr + zSqr)
					dist := math.Sqrt(mySqr)
					sum += dist
				}
			}

			avg := sum / 8

			if avg < 65536/QE {
				img.Set(x, y, color.NRGBA{
					R: uint8((255) & 255),
					G: uint8((255) << 1 & 255),
					B: uint8((255) << 2 & 255),
					A: 255,
				})
			} else {
				img.Set(x, y, color.NRGBA{
					R: uint8((0) & 255),
					G: uint8((0) << 1 & 255),
					B: uint8((0) << 2 & 255),
					A: 255,
				})
			}
		}
	}

	return img
}

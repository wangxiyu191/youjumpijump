package jump

import (
	"image"
	"math"
)

func edgeDetection(pic image.Image, QE float64) [][]int {
	limit := pic.Bounds().Size()

	img := [][]int{}
	for x := 0; x < limit.X; x++ {
		line := []int{}

		for y := 0; y < limit.Y; y++ {
			if (x < ExcludedWeight || x > limit.X-ExcludedWeight) ||
				(y < ExcludedHeight || y > limit.Y-ExcludedHeight) {
				line = append(line, 0)
				continue
			}

			sum := 0.0
			pixeloO := pic.At(x, y)
			x2, y2, z2, _ := pixeloO.RGBA()

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
				line = append(line, 0)
			} else {
				line = append(line, 1)
			}
		}
		img = append(img, line)
	}

	return img
}

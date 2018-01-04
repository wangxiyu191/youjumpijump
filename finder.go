package jump

import (
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/nfnt/resize"
)

const ExcludedHeight = 350 // 排除一些高和宽，避免干扰也减少计算量
const ExcludedWeight = 35

const BottleRadius = 20 // 跳瓶大小
const BottleHeight = 108

func findBlock(img [][]int) [4]int {
	axisX, axisY, weight := 0, 0, 0

	limitX, limitY := len(img), len(img[0])
	topY := 0
	// 找顶部图形中轴线
	for y := ExcludedHeight; y < limitY-ExcludedHeight; y++ {
		for x := ExcludedWeight; x < limitX-ExcludedWeight && axisX == 0; x++ {
			if img[x][y] == 1 {
				// 拒绝毛刺
				minX, maxX := 0, 0
				for x2 := ExcludedWeight; x2 < limitX-ExcludedWeight; x2++ {
					if img[x2][y+3] == 1 {
						if minX == 0 || x2 < minX {
							minX = x2
						}
						if maxX == 0 || x2 > maxX {
							maxX = x2
						}
					}
				}
				axisX = (minX + maxX) / 2
				topY = y
			}
		}
	}

	lastY := 0
	maxWeight := int(math.Min(float64(limitX-ExcludedWeight-axisX), float64(axisX-ExcludedWeight)))
	// 从中轴线找最大对称点
	for w := 1; w < maxWeight-1 && axisY == 0; w++ {
		line := []int{}
		for y := ExcludedHeight; y < limitY-ExcludedHeight-1; y++ {
			if img[axisX-w][y] == 1 && img[axisX+w][y] == 1 {
				line = append(line, y)
			}
		}
		for i := 0; i < len(line); i++ {
			if i == len(line)-1 || line[i] != line[i+1]-1 {
				avgY := (line[i]-line[0])/2 + line[0]
				// 期望排除方块旁边的竖线，也期望包含圆旁边的竖的像素点，圆旁边的竖线的像素点会少一些
				if lastY != 0 && avgY-lastY > 10 {
					axisY = lastY
					weight = w * 2
				}
				lastY = avgY
				break
			}
		}
		if len(line) == 0 {
			axisY = lastY
			weight = w * 2
		}
	}

	for x := axisX; x < axisX+weight/2; x++ {

	}

	return [4]int{axisX, axisY, weight, (axisY - topY) * 2}
}

func findBottle(img [][]int) [2]int {
	axisX, axisY := 0, 0

	limitX, limitY := len(img), len(img[0])
	for y := ExcludedHeight; y < limitY-ExcludedHeight; y++ {
		for x := ExcludedWeight; x < limitX-ExcludedWeight && axisX == 0; x++ {

			// 跳瓶的顶部圆
			minX, maxX := 0, 0
			minY, maxY := 0, 0
			for r := BottleRadius - 3; r < BottleRadius+3; r += 2 {
				if img[x+r][y] == 1 && img[x-r][y] == 1 && img[x][y+r] == 1 && img[x][y-r] == 1 {
					w := int(math.Pow(math.Pow(float64(r), 2)-math.Pow(float64(r/2), 2), 0.5))
					if img[x+r/2][y+w] == 1 && img[x+r/2][y-w] == 1 {
						w := int(math.Pow(math.Pow(float64(r), 2)-math.Pow(float64(r/4), 2), 0.5))
						if img[x+r/4][y+w] == 1 && img[x+r/4][y-w] == 1 {
							if minX == 0 || x < minX {
								minX = x
							}
							if maxX == 0 || x > maxX {
								maxX = x
							}
							if minY == 0 || y < minY {
								minY = y
							}
							if maxY == 0 || y > maxY {
								maxY = y
							}
						}
					}
				}
			}

			// 跳瓶重心
			axisX = (minX + maxX) / 2
			axisY = (minY+maxY)/2 + BottleHeight
		}
	}
	return [2]int{axisX, axisY}
}

func Find(pic image.Image) ([]int, []int) {
	start := time.Now()
	pic = resize.Resize(720, 0, pic, resize.Lanczos3)
	img := edgeDetection(pic, 100.0)
	savePNG("jump.720.png", pic)
	log.Printf("prepare took %s", time.Since(start))

	start = time.Now()
	log.Printf("%v %v", findBlock(img), findBottle(img))
	log.Printf("findBlock took %s", time.Since(start))

	f, _ := os.OpenFile("jump.edge.png", os.O_WRONLY|os.O_CREATE, 0600)
	png.Encode(f, toImg(img))
	f.Close()
	return nil, nil
}

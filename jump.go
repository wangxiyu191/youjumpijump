package jump

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

var basePath string

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath = filepath.Dir(ex)

	if ok, _ := Exists(basePath + "/debugger"); !ok {
		os.MkdirAll(basePath+"/debugger", os.ModePerm)
	}

	os.Remove(basePath + "/debugger/debug.log")
	logFile, _ := os.OpenFile(basePath+"/debugger/debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func Debugger() {
	if ok, _ := Exists(basePath + "/jump.png"); ok {
		os.Rename(basePath+"/jump.png", basePath+"/debugger/"+strconv.Itoa(TimeStamp())+".png")

		files, err := ioutil.ReadDir(basePath + "/debugger/")
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			fname := f.Name()
			ext := filepath.Ext(fname)
			name := fname[0 : len(fname)-len(ext)]
			if ts, err := strconv.Atoi(name); err == nil {
				if TimeStamp()-ts > 10 {
					os.Remove(basePath + "/debugger/" + fname)
				}
			}
		}
	}
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TimeStamp() int {
	return int(time.Now().UnixNano() / int64(time.Second))
}

func Distance(a, b []int) float64 {
	return math.Pow(math.Pow(float64(a[0]-b[0]), 2)+math.Pow(float64(a[1]-b[1]), 2), 0.5)
}

func getRGB(nowColor []int, src *image.RGBA, x int, y int) {
	offest := src.PixOffset(x, y)
	for i := 0; i < 3; i++ {
		nowColor[i] = int(src.Pix[offest+i])
	}
}

func setRGB(nowColor []int, src *image.RGBA, x int, y int) {
	offest := src.PixOffset(x, y)
	for i := 0; i < 3; i++ {
		src.Pix[offest+i] = uint8(nowColor[i])
	}
}

func colorSimilar(a, b []int, distance float64) bool {
	return (math.Abs(float64(a[0]-b[0])) < distance) && (math.Abs(float64(a[1]-b[1])) < distance) && (math.Abs(float64(a[2]-b[2])) < distance)
}

func Find(src_raw image.Image) ([]int, []int) {
	src_raw = resize.Resize(720, 0, src_raw, resize.NearestNeighbor)
	go func() {
		f, _ := os.OpenFile("jump.720.png", os.O_WRONLY|os.O_CREATE, 0600)
		png.Encode(f, src_raw)
		f.Close()
	}()

	bounds := src_raw.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	src := image.NewRGBA(src_raw.Bounds())

	switch src_raw.ColorModel() {
	case color.RGBAModel:
		src = src_raw.(*image.RGBA)
	case color.RGBA64Model:
		src.Pix = src_raw.(*image.RGBA64).Pix
		src.Stride = src_raw.(*image.RGBA64).Stride
	case color.NRGBAModel:
		src.Pix = src_raw.(*image.NRGBA).Pix
		src.Stride = src_raw.(*image.NRGBA).Stride
	case color.NRGBA64Model:
		src.Pix = src_raw.(*image.NRGBA64).Pix
		src.Stride = src_raw.(*image.NRGBA64).Stride
	}
	debugImg := image.NewRGBA(src_raw.Bounds())
	copy(debugImg.Pix, src.Pix)
	redColor := []int{255, 0, 0}
	blueColor := []int{0, 0, 255}
	jumpCubeColor := []int{54, 52, 92}
	points := [][]int{}

	nowColor := []int{0, 0, 0}

	bgColors := [][]int{}

	for y := 0; y < h; y++ {
		bgColor := []int{0, 0, 0}
		getRGB(bgColor, src, w-25, y)
		bgColors = append(bgColors, bgColor)
	}

	for y := 0; y < h; y++ {
		line := 0
		for x := 0; x < w; x++ {
			getRGB(nowColor, src, x, y)

			if colorSimilar(nowColor, jumpCubeColor, 20) {
				line++
				setRGB(bgColors[y], src, x, y)
				setRGB(redColor, debugImg, x, y)
			} else {
				if y > 350 && x-line > 10 && line > 30 {
					points = append(points, []int{x - line/2, y, line})
				}
				line = 0
			}
		}
	}
	jumpCube := []int{0, 0, 0}
	for _, point := range points {
		if point[2] > jumpCube[2] {
			jumpCube = point
		}
	}
	jumpCube = []int{jumpCube[0], jumpCube[1]}
	if jumpCube[0] == 0 {
		return nil, nil
	}

	// possible := [][]int{}

	minHeight := h
	minHeightX := -1

	for x := 0; x < w; x++ {
		for y := 350; y < h; y++ {
			bgColor := bgColors[y]
			getRGB(nowColor, src, x, y)
			if !colorSimilar(nowColor, bgColor, 5) {
				if y < minHeight {
					minHeight = y
					minHeightX = x
				}
			}

		}
	}

	target := []int{minHeightX, minHeight + 70}
	for x := target[0] - 5; x < target[0]+5; x++ {
		for y := target[1] - 5; y < target[1]+5; y++ {
			setRGB(blueColor, debugImg, x, y)
		}
	}

	for x := jumpCube[0] - 5; x < jumpCube[0]+5; x++ {
		for y := jumpCube[1] - 5; y < jumpCube[1]+5; y++ {
			setRGB(blueColor, debugImg, x, y)
		}
	}

	go func() {
		f, _ := os.OpenFile("jump.720.debug.png", os.O_WRONLY|os.O_CREATE, 0600)
		png.Encode(f, debugImg)
		f.Close()
	}()

	return jumpCube, target
}

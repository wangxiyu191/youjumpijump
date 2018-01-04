package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	jump "github.com/faceair/youjumpijump"
)

var similar *jump.Similar

func screenshot() image.Image {
	out, err := exec.Command("/system/bin/screencap").Output()
	if err != nil {
		panic("screenshot failed")
	}
	bound := image.Rect(0, 0, int(binary.LittleEndian.Uint32(out[0:])), int(binary.LittleEndian.Uint32(out[4:])))
	src := image.NewRGBA(bound)
	src.Pix = out[12:]
	return src
}

func main() {
	defer func() {
		jump.Debugger()
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("the program has crashed, press any key to exit")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var inputRatio float64
	var err error
	if len(os.Args) > 1 {
		inputRatio, err = strconv.ParseFloat(os.Args[1], 10)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Print("input jump ratio (recommend 2.04):")
		_, err = fmt.Scanln(&inputRatio)
		if err != nil {
			log.Printf("input is empty, will use 2.04 as default ratio")
			inputRatio = 2.04
		}
	}

	similar = jump.NewSimilar(inputRatio)

	for {
		jump.Debugger()

		src := screenshot()

		start, end := jump.Find(src)
		if start == nil {
			log.Print("can't find the starting point，please export the debugger directory")
			break
		} else if end == nil {
			log.Print("can't find the end point，please export the debugger directory")
			break
		}

		scale := float64(src.Bounds().Max.X) / 720
		nowDistance := jump.Distance(start, end)
		// TODO: https://github.com/faceair/youjumpijump/issues/61
		// similarDistance, nowRatio := similar.Find(nowDistance)
		similarDistance, nowRatio := 0.0, inputRatio

		log.Printf("from:%v to:%v distance:%.2f similar:%.2f ratio:%v press:%.2fms ", start, end, nowDistance, similarDistance, nowRatio, nowDistance*nowRatio)

		_, err = exec.Command("/system/bin/sh", "/system/bin/input", "swipe", strconv.FormatFloat(float64(start[0])*scale, 'f', 0, 32), strconv.FormatFloat(float64(start[1])*scale, 'f', 0, 32), strconv.FormatFloat(float64(end[0])*scale, 'f', 0, 32), strconv.FormatFloat(float64(end[1])*scale, 'f', 0, 32), strconv.Itoa(int(nowDistance*nowRatio))).Output()
		if err != nil {
			panic("touch failed")
		}

		go func() {
			// TODO: https://github.com/faceair/youjumpijump/issues/61
			return
			time.Sleep(time.Millisecond * 170)
			src := screenshot()

			finally, _ := jump.Find(src)
			if finally != nil {
				finallyDistance := jump.Distance(start, finally)
				finallyRatio := (nowDistance * nowRatio) / finallyDistance

				if finallyRatio > nowRatio/2 && finallyRatio < nowRatio*2 {
					similar.Add(finallyDistance, finallyRatio)
				}
			}
		}()

		time.Sleep(time.Millisecond * 1500)
	}
}

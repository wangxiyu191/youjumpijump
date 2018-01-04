package main

import (
	"image/png"
	"os"

	"github.com/faceair/youjumpijump"
)

func main() {
	inFile, err := os.Open("jump.png")
	if err != nil {
		panic(err)
	}
	src, err := png.Decode(inFile)
	if err != nil {
		panic(err)
	}
	inFile.Close()

	jump.Find(src)
}

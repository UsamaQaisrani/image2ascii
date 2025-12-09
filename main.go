package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

func main() {
	imagePath := "chess.jpg"
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error while loading image file: ", err)
		os.Exit(1)
	}
	pixels, err := getImagePixels(file)
	if err != nil {
		log.Fatal("Error while getting pixel array:", err)
		os.Exit(1)
	}
	fmt.Println(pixels)
}

func getImagePixels(file io.ReadCloser) ([][]Pixel, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}
	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

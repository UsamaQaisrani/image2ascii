package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

var characters = [10]string{" ", ".", ":", "-", "=", "+", "*", "#", "%", "@"}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please add path to the image: go run . <path>")
	}
	imagePath := args[0]
	fmt.Println(imagePath)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error while loading image file: ", err)
		os.Exit(1)
	}
	height, width, pixels, err := getImageInfo(file)
	if err != nil {
		log.Fatal("Error while getting pixel array:", err)
		os.Exit(1)
	}
	asciiPixels := buildAscii(pixels, height, width)
	for _, row := range asciiPixels {
		line := ""
		for _, ch := range row {
			line += ch
		}
		fmt.Println(line)
	}
}

func resizePixels(pixels [][]Pixel, originalHeight, originalWidth, newHeight, newWidth int) [][]Pixel {
	resizedPixels := make([][]Pixel, newHeight)
	for i := range resizedPixels {
		resizedPixels[i] = make([]Pixel, newWidth)
	}
	for i := 0; i < newHeight; i++ {
		for j := 0; j < newWidth; j++ {
			originalY := int(float64(i) * float64(originalHeight) / float64(newHeight))
			originalX := int(float64(j) * float64(originalWidth) / float64(newWidth))
			if originalY >= originalHeight {
				originalY = originalHeight - 1
			}
			if originalX >= originalWidth {
				originalX = originalWidth - 1
			}
			resizedPixels[i][j] = pixels[originalY][originalX]
		}
	}
	return resizedPixels
}

func buildAscii(pixels [][]Pixel, height, width int) [][]string {
	newWidth := 60
	verticalCorrection := 0.5
	newHeight := int(float64(height) * (float64(newWidth) / float64(width)) * verticalCorrection)
	resizedPixels := resizePixels(pixels, height, width, newHeight, newWidth)
	asciiPixels := make([][]string, newHeight)
	for i, row := range resizedPixels {
		asciiPixels[i] = make([]string, newWidth)
		for j, col := range row {
			intensity := getIntensityOfPixel(col)
			index := getCharacterIndex(intensity)
			asciiPixels[i][j] = characters[index]
		}
	}
	return asciiPixels
}

func getImageInfo(file io.ReadCloser) (height int, width int, pixelsList [][]Pixel, err error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, nil, err
	}
	bounds := img.Bounds()
	width, height = bounds.Max.X, bounds.Max.Y
	pixels := make([][]Pixel, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]Pixel, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = rgbaToPixel(img.At(x, y).RGBA())
		}
	}
	return height, width, pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func getIntensityOfPixel(pixel Pixel) float64 {
	return 0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B)
}

func getCharacterIndex(intensity float64) int {
	index := int(math.Floor(intensity / (256.0 / 10.0)))
	if index >= 10 {
		index = 9
	}
	return index
}

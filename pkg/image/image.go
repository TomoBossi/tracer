package image

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"tracer/pkg/utils"
)

type Threshold interface {
	eval(pixels [][]color.Gray, x, y int) bool
}

type ThresholdAbsolute struct {
	LessOrEqual uint8
}

func (t ThresholdAbsolute) eval(pixels [][]color.Gray, x, y int) bool {
	return pixels[y][x].Y <= t.LessOrEqual
}

type ThresholdRelativeArea struct {
	LessOrEqual uint8
	GreaterDiff uint8
	Radius      int
}

func (t ThresholdRelativeArea) eval(pixels [][]color.Gray, x, y int) bool {
	average := averageGray(subPixelsCirc(pixels, x, y, t.Radius))
	return uint8(utils.Abs(int(average.Y)-int(pixels[y][x].Y))) >= t.GreaterDiff && pixels[y][x].Y <= t.LessOrEqual
}

func CreateBinaryImage(binary [][]bool) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, len(binary[0]), len(binary)))
	var pixel uint8 = 0
	for y, row := range binary {
		for x, is_traced := range row {
			if is_traced {
				pixel = 0
			} else {
				pixel = 255
			}
			img.Set(x, y, color.Gray{pixel})
		}
	}
	return img
}

func BinaryPixels(pixels [][]color.Gray, threshold Threshold) [][]bool {
	binary := make([][]bool, len(pixels))
	for y := range pixels {
		binary[y] = make([]bool, len(pixels[y]))
		for x := range pixels[y] {
			binary[y][x] = threshold.eval(pixels, x, y)
		}
	}
	return binary
}

func flatten(pixels [][]color.Gray) []color.Gray {
	flattenedPixels := make([]color.Gray, count(pixels))
	i := 0
	for _, row := range pixels {
		for _, pixel := range row {
			flattenedPixels[i] = pixel
			i++
		}
	}
	return flattenedPixels
}

func averageGray(pixels [][]color.Gray) color.Gray {
	sum := 0.
	for _, row := range pixels {
		for _, pixel := range row {
			sum += float64(pixel.Y)
		}
	}
	return color.Gray{uint8(sum / float64(count(pixels)))}
}

func count(pixels [][]color.Gray) int {
	q := 0
	for i := range pixels {
		q += len(pixels[i])
	}
	return q
}

func rgbaToGray(pixel color.RGBA) color.Gray {
	return color.Gray{uint8(0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B))}
}

func RgbaToGrayPixels(pixels [][]color.RGBA) [][]color.Gray {
	grayscalePixels := make([][]color.Gray, len(pixels))
	for y, row := range pixels {
		grayscalePixels[y] = make([]color.Gray, len(row))
		for x, pixel := range row {
			grayscalePixels[y][x] = rgbaToGray(pixel)
		}
	}
	return grayscalePixels
}

func subPixelsCirc[T color.RGBA | color.Gray](pixels [][]T, x, y, radius int) [][]T {
	rowFrom, rowTo := max(0, y-radius), min(len(pixels), y+radius)
	colFrom, colTo := max(0, x-radius), min(len(pixels[0]), x+radius)
	sub := make([][]T, 0, rowTo-rowFrom)
	for i, row := range pixels[rowFrom:rowTo] {
		subRow := make([]T, 0, colTo-colFrom)
		for j, pixel := range row[colFrom:colTo] {
			if int(utils.L2(j, i, x-colFrom, y-rowFrom)) <= radius {
				subRow = append(subRow, pixel)
			}
		}
		sub = append(sub, subRow)
	}
	return sub
}

func subPixelsRect[T color.RGBA | color.Gray](pixels [][]T, rowFrom, rowTo, colFrom, colTo int) [][]T {
	if rowTo == 0 {
		rowTo = len(pixels)
	}

	if colTo == 0 && rowTo > 0 {
		colTo = len(pixels[0])
	}

	sub := make([][]T, rowTo-rowFrom)
	for i := rowFrom; i < rowTo; i++ {
		sub[i-rowFrom] = pixels[i][colFrom:colTo]
	}
	return sub
}

func LoadRgbaPixels(path string) [][]color.RGBA {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("LoadRgbaPixels: Error opening file -", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("LoadRgbaPixels: Error decoding image -", err)
	}
	if img == nil {
		fmt.Println("LoadRgbaPixels: Decoded image is nil")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := make([][]color.RGBA, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]color.RGBA, width)
		for x := 0; x < width; x++ {
			col := img.At(x, y)
			r, g, b, a := col.RGBA()
			rgba := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			pixels[y][x] = rgba
		}
	}
	return pixels
}

func CreateGrayscaleImage(pixels [][]color.Gray) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, len(pixels[0]), len(pixels)))
	for y, row := range pixels {
		for x, pixel := range row {
			img.Set(x, y, pixel)
		}
	}
	return img
}

func CreateRgbaImage(pixels [][]color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, len(pixels[0]), len(pixels)))
	for y, row := range pixels {
		for x, pixel := range row {
			img.Set(x, y, pixel)
		}
	}
	return img
}

func SavePngImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("SavePngImage: Error creating file - %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("SavePngImage: Error encoding image - %v", err)
	}

	fmt.Println("SavePngImage: Image saved as", path)
	return nil
}

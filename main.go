package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

type Threshold interface {
	eval(pixels [][]color.Gray, x, y int) bool
}

type ThresholdAbsolute struct {
	less_or_equal uint8
}

func (t ThresholdAbsolute) eval(pixels [][]color.Gray, x, y int) bool {
	return pixels[y][x].Y <= t.less_or_equal
}

type ThresholdRelativeArea struct {
	less_or_equal uint8
	greater_diff  uint8
	radius        int
}

func (t ThresholdRelativeArea) eval(pixels [][]color.Gray, x, y int) bool {
	average := average_gray(sub_pixels_circ(pixels, x, y, t.radius))
	return uint8(abs(int(average.Y)-int(pixels[y][x].Y))) >= t.greater_diff && pixels[y][x].Y <= t.less_or_equal
}

func main() {
	pixels := rgba_to_gray_pixels(load_rgba_pixels("/home/tomo/Downloads/lain.jpg"))
	binary := binary_pixels(pixels, ThresholdRelativeArea{135, 20, 10})
	img := create_binary_image(binary)
	save_png_image(img, "/home/tomo/Downloads/test.png")
}

func abs(n int) int {
	if n < 0 {
		return int(-n)
	}
	return int(n)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func create_binary_image(binary [][]bool) *image.Gray {
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

func binary_pixels(pixels [][]color.Gray, threshold Threshold) [][]bool {
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
	flattened_pixels := make([]color.Gray, count(pixels))
	i := 0
	for _, row := range pixels {
		for _, pixel := range row {
			flattened_pixels[i] = pixel
			i++
		}
	}
	return flattened_pixels
}

func average_gray(pixels [][]color.Gray) color.Gray {
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

func rgba_to_gray(pixel color.RGBA) color.Gray {
	return color.Gray{uint8(0.299*float64(pixel.R) + 0.587*float64(pixel.G) + 0.114*float64(pixel.B))}
}

func rgba_to_gray_pixels(pixels [][]color.RGBA) [][]color.Gray {
	grayscale_pixels := make([][]color.Gray, len(pixels))
	for y, row := range pixels {
		grayscale_pixels[y] = make([]color.Gray, len(row))
		for x, pixel := range row {
			grayscale_pixels[y][x] = rgba_to_gray(pixel)
		}
	}
	return grayscale_pixels
}

func l2(x1, y1, x2, y2 int) float64 {
	dx, dy := x1-x2, y1-y2
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func sub_pixels_circ[T color.RGBA | color.Gray](pixels [][]T, x, y, radius int) [][]T {
	row_from, row_to := max(0, y-radius), min(len(pixels), y+radius)
	col_from, col_to := max(0, x-radius), min(len(pixels[0]), x+radius)
	sub_rect := sub_pixels_rect(pixels, row_from, row_to, col_from, col_to)
	center_x, center_y := y-row_from, x-col_from
	sub := make([][]T, 0, row_to-row_from)
	for y, row := range sub_rect {
		sub_row := make([]T, 0, col_to-col_from)
		for x, pixel := range row {
			if int(l2(x, y, center_x, center_y)) <= radius {
				sub_row = append(sub_row, pixel)
			}
		}
		sub = append(sub, sub_row)
	}
	return sub
}

func sub_pixels_rect[T color.RGBA | color.Gray](pixels [][]T, row_from, row_to, col_from, col_to int) [][]T {
	if row_to == 0 {
		row_to = len(pixels)
	}

	if col_to == 0 && row_to > 0 {
		col_to = len(pixels[0])
	}

	sub := make([][]T, row_to-row_from)
	for i := row_from; i < row_to; i++ {
		sub[i-row_from] = pixels[i][col_from:col_to]
	}
	return sub
}

func load_rgba_pixels(path string) [][]color.RGBA {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("load_rgba_pixels: Error opening file -", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("load_rgba_pixels: Error decoding image -", err)
	}
	if img == nil {
		fmt.Println("load_rgba_pixels: Decoded image is nil")
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

func create_grayscale_image(pixels [][]color.Gray) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, len(pixels[0]), len(pixels)))
	for y, row := range pixels {
		for x, pixel := range row {
			img.Set(x, y, pixel)
		}
	}
	return img
}

func create_rgba_image(pixels [][]color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, len(pixels[0]), len(pixels)))
	for y, row := range pixels {
		for x, pixel := range row {
			img.Set(x, y, pixel)
		}
	}
	return img
}

func save_png_image(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save_png_image: Error creating file - %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("save_png_image: Error encoding image - %v", err)
	}

	fmt.Println("save_png_image: Image saved as", path)
	return nil
}

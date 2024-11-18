package main

import "tracer/pkg/image"

func main() {
	pixels := image.RgbaToGrayPixels(image.LoadRgbaPixels("/home/tomo/Downloads/lain.jpg"))
	binary := image.BinaryPixels(pixels, image.ThresholdRelativeArea{LessOrEqual: 135, GreaterDiff: 20, Radius: 13})
	img := image.CreateBinaryImage(binary)
	image.SavePngImage(img, "/home/tomo/Downloads/test.png")
}

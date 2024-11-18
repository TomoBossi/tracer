package main

import (
	"fmt"
	"tracer/pkg/cluster"
	"tracer/pkg/image"
)

func main() {
	pixels := image.RgbaToGrayPixels(image.LoadRgbaPixels("./assets/clusterize_medium.png"))
	binary := image.BinaryPixels(pixels, image.ThresholdAbsolute{LessOrEqual: 150})
	//fmt.Println(len(cluster.Clusterize(binary, 5)))
	fmt.Println(cluster.Clusterize(binary, 5))
	// img := image.CreateBinaryImage(binary)
	// image.SavePngImage(img, "./out/test.png")
}

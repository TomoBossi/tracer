package cluster

type Point struct {
	X, Y int
}

func Clusterize(binary [][]bool, maxDistance int) [][]Point {
	clusters := make([][]Point, 0)
	copy := deepCopy(binary)
	for y := range copy {
		for x := range copy[y] {
			if copy[y][x] {
				cluster := make([]Point, 0, 1)
				neighbors := []Point{{x, y}}
				for len(neighbors) > 0 {
					for _, n := range neighbors {
						copy[n.Y][n.X] = false
					}
					n := neighbors[0]
					nX, nY := n.X, n.Y
					cluster = append(cluster, n)
					neighbors = append(neighbors[1:], scan(copy, nX, nY, maxDistance)...)
				}
				clusters = append(clusters, cluster)
			}
		}
	}
	return clusters
}

func deepCopy(binary [][]bool) [][]bool {
	copy := make([][]bool, len(binary))
	for y, row := range binary {
		copy[y] = make([]bool, len(binary[y]))
		for x, value := range row {
			copy[y][x] = value
		}
	}
	return copy
}

func scan(binary [][]bool, x, y, maxDistance int) []Point {
	neighbors := make([]Point, 0, 4*maxDistance*maxDistance)
	width, height := len(binary[0]), len(binary)
	for d := 1; d <= maxDistance; d++ {
		nX, nY := x+d, y+d
		for i := 0; i < d*2+1; i++ {
			if inRange(width, height, nX-d*2, nY-i) && binary[nY-i][nX-d*2] {
				neighbors = append(neighbors, Point{nX - d*2, nY - i})
			}
			if inRange(width, height, nX, nY-i) && binary[nY-i][nX] {
				neighbors = append(neighbors, Point{nX, nY - i})
			}
			if i != 0 && i != d*2 {
				if inRange(width, height, nX-i, nY) && binary[nY][nX-i] {
					neighbors = append(neighbors, Point{nX - i, nY})
				}
				if inRange(width, height, nX-i, nY-d*2) && binary[nY-d*2][nX-i] {
					neighbors = append(neighbors, Point{nX - i, nY - d*2})
				}
			}
		}
	}
	return neighbors
}

func inRange(width, height, x, y int) bool {
	return x > 0 && y > 0 && x < width && y < height
}

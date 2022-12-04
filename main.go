package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

var (
	pointsCount = flag.Int("p", 5, "define the amount of origin points (p > 0)")
	width       = flag.Int("w", 640, "define a custom width for the canvas")
	height      = flag.Int("h", 480, "define a custom height for the canvas")
	help        = flag.Bool("help", false, "prints a guid on how to use the program")
)
var originPoints []OriginPoint

type OriginPoint struct {
	X, Y  int
	Color color.Color
}

type PointDistance struct {
	Point    OriginPoint
	Distance float64
}

func init() {
	flag.Parse()

	if *help {
		fmt.Println("Welcome to the random image thing (whatever it is called) generator!")
		fmt.Println("How to use:")
		fmt.Println("\tIf you only want an image:")
		fmt.Printf("\tjust run .\\graphic.exe and it will generate a %d by %d image with %d points.\n\n", *width, *height, *pointsCount)
		fmt.Println("If you want to customize your image size and amount of points, check the flag description:")
		flag.PrintDefaults()
	}

	if *pointsCount < 0 {
		panic("cannot have a negative value for points")
	}

	rand.Seed(time.Now().UnixMilli())

	originPoints = make([]OriginPoint, *pointsCount)
	for i := 0; i < *pointsCount; i++ {
		originPoints[i] = OriginPoint{
			X: rand.Intn(*width),
			Y: rand.Intn(*height),
			Color: color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 255,
			},
		}
	}
}

func selectColor(array []PointDistance) color.Color {
	min := array[0].Distance
	color := array[0].Point.Color

	for _, value := range array {
		if min > value.Distance {
			min = value.Distance
			color = value.Point.Color
		}
	}
	return color
}

func distanceFromAllPoints(x, y int) color.Color {
	points := make([]PointDistance, len(originPoints))

	for i, point := range originPoints {
		xs := math.Pow(float64(point.X-x), 2)
		ys := math.Pow(float64(point.Y-y), 2)

		d := math.Sqrt(xs + ys)

		points[i] = PointDistance{point, d}
	}

	return selectColor(points)
}

func main() {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{*width, *height},
	})

	for x := 0; x < *width; x++ {
		for y := 0; y < *height; y++ {
			c := distanceFromAllPoints(x, y)
			img.Set(x, y, c)
		}
	}

	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

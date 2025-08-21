package structs

import "math"

type Rectangle struct {
	width float64
	height float64
}

type Circle struct {
	radius float64
}

func (d *Circle) Area() float64 {
	return math.Pi 
}

func (d *Rectangle) Area() float64 {
	return d.height * d.width
}
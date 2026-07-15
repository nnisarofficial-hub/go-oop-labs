package main

import (
	"fmt"
	"math"
)

// func main() {}

type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2) // Circle Radius Formula "A = pi x radius^2"
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius //Perimeter (Circumference) Formula "C = 2 x pi x radius"
}

func (c Circle) Name() string { return "Circle" }

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Height * r.Width // Rectangle Formula "Area = height x width"
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width) //Perimeter (Circumference) Formula "Rectangle = 2 x (height x width"
}

func (r Rectangle) Name() string { return "Rectangle" }

type Triangle struct {
	A, B, C float64
}

func (t Triangle) Area() float64 {
	semiPerimeter := t.Perimeter() / 2
	return math.Sqrt(semiPerimeter * (semiPerimeter - t.A) * (semiPerimeter - t.B) * (semiPerimeter - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) Name() string { return "Triangle" }

func PrintShapeInfo(s Shape) {
	fmt.Printf("%s\n", s.Name())
	fmt.Printf("  %-12s%5.2f\n", "Area:", s.Area())
	fmt.Printf("  %-12s%5.2f\n", "Perimeter:", s.Perimeter())
}

func main() {
	circle := Circle{Radius: 10}
	rectangle := Rectangle{Width: 15, Height: 20}
	triangle := Triangle{A: 6, B: 5, C: 9}

	PrintShapeInfo(circle)
	PrintShapeInfo(rectangle)
	PrintShapeInfo(triangle)

	shapes := []Shape{circle, rectangle, triangle}
	bigArea := LargestShape(shapes)
	if bigArea != nil {
		fmt.Printf("\nLargest shape by area: %s (%.2f)\n", bigArea.Name(), bigArea.Area())
	}
}

func LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}
	largest := shapes[0]
	for _, shape := range shapes {
		if shape.Area() > largest.Area() {
			largest = shape
			// fmt.Print(shape)
		}
	}
	return largest
}

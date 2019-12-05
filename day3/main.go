package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Abs returns the absolute value for the input
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Intersection is a container for an intersection point and
// the number of steps taken to the point
type Intersection struct {
	Point Point
	Steps int
}

// Point is a representation of a coordinate point
type Point struct {
	X int
	Y int
}

// DistanceToPoint calculates the `Manhatten Distance` from
// one Point to another Point
func (p Point) DistanceToPoint(p2 Point) int {
	return Abs(p.X-p2.X) + Abs(p.Y-p2.Y)
}

// Path is a container for an array of Point(s) created by a
// series of Vector(s)
type Path struct {
	Points []Point
}

// Intersections returns an array of all intersections between
// two paths
func (p Path) Intersections(p2 Path) []Intersection {
	intersect := []Intersection{}

	for i, pointA := range p.Points {
		for j, pointB := range p2.Points {
			if pointA.X == pointB.X && pointA.Y == pointB.Y {
				in := Intersection{
					Point: Point{X: pointA.X, Y: pointB.Y},
					Steps: i + 1 + j + 1}

				intersect = append(intersect, in)
			}
		}
	}

	return intersect
}

// NewPath creates a Path from an array of Vector(s)
func NewPath(v []Vector) Path {
	p := Path{
		Points: []Point{}}

	x := 0
	y := 0

	for _, vec := range v {
		switch vec.Direction {
		case "R":
			for i := x + 1; i <= x+vec.Magnitude; i++ {
				point := Point{
					X: i,
					Y: y}
				p.Points = append(p.Points, point)
			}
			x += vec.Magnitude
		case "L":
			for i := x - 1; i >= x-vec.Magnitude; i-- {
				point := Point{
					X: i,
					Y: y}
				p.Points = append(p.Points, point)
			}
			x -= vec.Magnitude
		case "U":
			for i := y + 1; i <= y+vec.Magnitude; i++ {
				point := Point{
					X: x,
					Y: i}
				p.Points = append(p.Points, point)
			}
			y += vec.Magnitude
		case "D":
			for i := y - 1; i >= y-vec.Magnitude; i-- {
				point := Point{
					X: x,
					Y: i}
				p.Points = append(p.Points, point)
			}
			y -= vec.Magnitude
		}
	}

	return p
}

// Vector is a representation of a vector based on the spec
// [A-Z][0-9]+
type Vector struct {
	Direction string
	Magnitude int
}

// NewVector creates a vector from a string
func NewVector(s string) Vector {
	v := Vector{
		Direction: "X",
		Magnitude: 0}

	r := []rune(s)

	v.Direction = string(r[0])
	m := string(r[1:])

	mag, err := strconv.Atoi(m)
	if err != nil {
		panic("parse error")
	}
	v.Magnitude = mag

	return v
}

// NewVectorList creates an array of Vector(s) from an array
// of strings
func NewVectorList(s []string) []Vector {
	vectorLst := []Vector{}

	for _, vStr := range s {
		vectorLst = append(vectorLst, NewVector(vStr))
	}

	return vectorLst
}

// ClosestPoint returns the closest calculated intersection
// from an origin point
func ClosestPoint(intersections []Intersection) int {
	origin := Point{X: 0, Y: 0}
	d := int(^uint(0) >> 1)
	for _, inter := range intersections {
		r := origin.DistanceToPoint(inter.Point)
		if r < d {
			d = r
		}
	}

	return d
}

// FewestSteps returns the least number of steps taken to get to
// an intersection
func FewestSteps(intersections []Intersection) int {
	d := int(^uint(0) >> 1)
	for _, inter := range intersections {
		if inter.Steps < d {
			d = inter.Steps
		}
	}

	return d
}

func main() {
	vec := [][]string{}
	file, err := os.Open("./vectors.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		va := []string{}
		values := strings.Split(scanner.Text(), ",")
		for _, v := range values {
			va = append(va, v)
		}
		vec = append(vec, va)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	path1 := NewPath(NewVectorList(vec[0]))
	path2 := NewPath(NewVectorList(vec[1]))

	intersections := path1.Intersections(path2)
	d := ClosestPoint(intersections)
	s := FewestSteps(intersections)

	println(fmt.Sprintf("part 1 - closest point: %v", d))
	println(fmt.Sprintf("part 2 - fewest steps: %v", s))
}

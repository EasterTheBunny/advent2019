package main

import (
	"testing"
)

func TestAbs(t *testing.T) {
	values := [][]int{[]int{-1, 1}, []int{1, 1}, []int{-23476, 23476}}

	for _, value := range values {
		x := Abs(value[0])
		if x != value[1] {
			t.Errorf("incorrect absolute value %v; expected %v", x, value[1])
		}
	}
}

func TestDistanceToPoint(t *testing.T) {
	args := []int{0, 1, 1, 2, -5, 2, -6, 29}
	expect := []int{2, 7, 6, 1, 32}

	for i := 0; i < len(args)-4; i++ {
		point1 := Point{X: args[i], Y: args[i+1]}
		point2 := Point{X: args[i+2], Y: args[i+3]}

		expected := expect[i]
		calc := point1.DistanceToPoint(point2)

		if calc != expected {
			t.Errorf("incorrect point distance %v; expected %v", calc, expected)
		}
	}
}

func TestNewVector(t *testing.T) {
	inp := []string{"U32", "L23", "D23422", "R9"}
	exp := []Vector{
		Vector{Direction: "U", Magnitude: 32},
		Vector{Direction: "L", Magnitude: 23},
		Vector{Direction: "D", Magnitude: 23422},
		Vector{Direction: "R", Magnitude: 9}}

	for i, value := range inp {
		v := NewVector(value)
		c := exp[i]

		if v.Direction != c.Direction {
			t.Errorf("incorrect vector direction %s; expected %s", v.Direction, c.Direction)
		}

		if v.Magnitude != c.Magnitude {
			t.Errorf("incorrect vector magnitude %v; expected %v", v.Magnitude, c.Magnitude)
		}
	}
}

func TestNewPath(t *testing.T) {
	vecStr := []string{"R8", "U5", "L5", "D3"}
	path := NewPath(NewVectorList(vecStr))

	if len(path.Points) != 21 {
		t.Errorf("path: %v", path.Points)
		t.Errorf("incorrect points length %v; expected %v", len(path.Points), 21)
	}

	expected := []Point{
		Point{X: 1, Y: 0},
		Point{X: 2, Y: 0},
		Point{X: 3, Y: 0},
		Point{X: 4, Y: 0},
		Point{X: 5, Y: 0},
		Point{X: 6, Y: 0},
		Point{X: 7, Y: 0},
		Point{X: 8, Y: 0},
		Point{X: 8, Y: 1},
		Point{X: 8, Y: 2},
		Point{X: 8, Y: 3},
		Point{X: 8, Y: 4},
		Point{X: 8, Y: 5},
		Point{X: 7, Y: 5},
		Point{X: 6, Y: 5},
		Point{X: 5, Y: 5},
		Point{X: 4, Y: 5},
		Point{X: 3, Y: 5},
		Point{X: 3, Y: 4},
		Point{X: 3, Y: 3},
		Point{X: 3, Y: 2}}

	for i, p := range expected {
		if p.X != path.Points[i].X {
			t.Errorf("incorrect x value %v for index %v; expected %v", path.Points[i].X, i, p.X)
		}

		if p.Y != path.Points[i].Y {
			t.Errorf("incorrect y value %v for index %v; expected %v", path.Points[i].Y, i, p.Y)
		}
	}
}

func TestClosestPoint(t *testing.T) {
	inter := []Intersection{
		Intersection{Point: Point{X: 1, Y: 8}, Steps: 5},
		Intersection{Point: Point{X: 1, Y: 5}, Steps: 8},
		Intersection{Point: Point{X: -3, Y: -1}, Steps: 23}}

	value := ClosestPoint(inter)
	expect := 4

	if value != expect {
		t.Errorf("incorrect closest point distance %v; expected %v", value, expect)
	}
}

func TestFewestSteps(t *testing.T) {
	inter := []Intersection{
		Intersection{Point: Point{X: 1, Y: 8}, Steps: 5},
		Intersection{Point: Point{X: 1, Y: 5}, Steps: 8},
		Intersection{Point: Point{X: -3, Y: -1}, Steps: 23}}

	value := FewestSteps(inter)
	expect := 5

	if value != expect {
		t.Errorf("incorrect fewest steps %v; expected %v", value, expect)
	}
}

func TestFullIntegration(t *testing.T) {
	vec1 := []string{"R75", "D30", "R83", "U83", "L12", "D49", "R71", "U7", "L72"}
	vec2 := []string{"U62", "R66", "U55", "R34", "D71", "R55", "D58", "R83"}

	path1 := NewPath(NewVectorList(vec1))
	path2 := NewPath(NewVectorList(vec2))

	intersections := path1.Intersections(path2)

	if len(intersections) == 0 {
		t.Error("no intersections found")
	}

	d := ClosestPoint(intersections)
	s := FewestSteps(intersections)

	if d != 159 {
		t.Errorf("incorrect lowest distance %v; expected %v", d, 159)
	}

	if s != 610 {
		t.Errorf("incorrect value for fewst steps %v; expected %v", s, 610)
	}
}

package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// Point ...
type Point struct {
	X   int
	Y   int
	Obj *Asteroid
}

// Asteroid ...
type Asteroid struct {
	InView []*Asteroid
}

// PointsBetween ...
func PointsBetween(a Point, b Point) []Point {
	p := []Point{}

	if (a.X > b.X && a.Y > b.Y) || (a.X < b.X && a.Y < b.Y) { // positive
		if a.X > b.X {
			a, b = b, a
		}
		a1 := float64(a.X)
		a2 := float64(a.Y)
		b1 := float64(b.X)
		b2 := float64(b.Y)

		slope := (b2 - a2) / (b1 - a1)

		for y := a.Y; y <= b.Y; y++ {
			for x := a.X; x <= b.X; x++ {
				if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
					continue
				}

				c1 := float64(x)
				c2 := float64(y)

				if c2-b2 == slope*(c1-b1) {
					p = append(p, Point{X: x, Y: y})
				}
			}
		}
	} else if a.X > b.X && a.Y < b.Y { // negative slope reversed points
		a1 := float64(a.X)
		a2 := float64(a.Y)
		b1 := float64(b.X)
		b2 := float64(b.Y)

		slope := (b2 - a2) / (b1 - a1)

		for y := a.Y; y <= b.Y; y++ {
			for x := a.X; x >= b.X; x-- {
				if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
					continue
				}

				c1 := float64(x)
				c2 := float64(y)

				if c2-b2 == slope*(c1-b1) {
					p = append(p, Point{X: x, Y: y})
				}
			}
		}
	} else if a.X < b.X && a.Y > b.Y { // negative slope
		a1 := float64(a.X)
		a2 := float64(a.Y)
		b1 := float64(b.X)
		b2 := float64(b.Y)

		slope := (b2 - a2) / (b1 - a1)

		for y := a.Y; y >= b.Y; y-- {
			for x := a.X; x <= b.X; x++ {
				if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
					continue
				}

				c1 := float64(x)
				c2 := float64(y)

				if c2-b2 == slope*(c1-b1) {
					p = append(p, Point{X: x, Y: y})
				}
			}
		}
	} else if a.X == b.X && a.Y != b.Y { // infinity slope
		if a.Y > b.Y {
			a, b = b, a
		}

		x := a.X
		for y := a.Y; y <= b.Y; y++ {
			if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
				continue
			}

			p = append(p, Point{X: x, Y: y})
		}
	} else if a.X != b.X && a.Y == b.Y { // zero slope
		if a.X > b.X {
			a, b = b, a
		}

		y := a.Y
		for x := a.X; x <= b.X; x++ {
			if (x == a.X && y == a.Y) || (x == b.X && y == b.Y) {
				continue
			}

			p = append(p, Point{X: x, Y: y})
		}
	} else { // points are the same
		return p
	}

	return p
}

// MakeGrid ...
func MakeGrid(r io.Reader) ([][]*Point, error) {
	grid := [][]*Point{}
	row := []*Point{}
	reader := bufio.NewReader(r)
	x := 0
	y := 0

	for {
		n, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		ch := string(n)

		point := Point{X: x, Y: y}
		x++
		switch ch {
		case ".":
			row = append(row, &point)
			continue
		case "#":
			point.Obj = &Asteroid{InView: []*Asteroid{}}
			row = append(row, &point)
		default:
			x = 0
			y++
			grid = append(grid, row)
			row = []*Point{}
		}
	}

	if len(row) > 0 {
		grid = append(grid, row)
	}

	return grid, nil
}

// ScanViews ...
func ScanViews(grid [][]*Point, p *Point) int {
	if p.Obj == nil {
		return 0
	}

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			point := grid[y][x]

			if point != nil && point.Obj != nil && !(point.X == p.X && point.Y == p.Y) {
				b := PointsBetween(*p, *point)
				canView := true
				for _, tp := range b {
					rp := grid[tp.Y][tp.X]
					if rp.Obj != nil {
						canView = false
						break
					}
				}

				if canView {
					p.Obj.InView = append(p.Obj.InView, point.Obj)
				}
			}
		}
	}

	return len(p.Obj.InView)
}

// BestView ...
func BestView(grid [][]*Point) *Point {
	var vantage *Point
	max := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			inView := ScanViews(grid, grid[y][x])
			if inView > max {
				vantage = grid[y][x]
				max = inView
			}
		}
	}

	return vantage
}

func main() {
	file, err := os.Open("./grid.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	grid, err := MakeGrid(file)
	if err != nil {
		panic("unexpected error")
	}

	max := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			inView := ScanViews(grid, grid[y][x])
			if inView > max {
				max = inView
			}
		}
	}
	println(max)

	// wrong answer: 291
}

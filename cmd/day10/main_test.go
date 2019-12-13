package main

import (
	"strings"
	"testing"
)

type PointPair struct {
	a Point
	b Point
	l int
}

func TestPointsBetween(t *testing.T) {
	inputs := []PointPair{
		PointPair{a: Point{X: 0, Y: 1}, b: Point{X: 4, Y: 3}, l: 1},
		PointPair{a: Point{X: 0, Y: 0}, b: Point{X: 9, Y: 9}, l: 8}}

	for _, input := range inputs {
		result := PointsBetween(input.a, input.b)
		if len(result) != input.l {
			t.Errorf("incorrect number of points between: %v; expected %v", len(result), input.l)
		}
	}
}

func TestMakeGrid(t *testing.T) {
	blank := &Asteroid{InView: []*Asteroid{}}
	expected := [][]*Point{
		[]*Point{&Point{X: 0, Y: 0}, &Point{X: 1, Y: 0, Obj: blank}, &Point{X: 2, Y: 0}},
		[]*Point{&Point{X: 0, Y: 1, Obj: blank}, &Point{X: 1, Y: 1}, &Point{X: 2, Y: 1, Obj: blank}},
		[]*Point{&Point{X: 0, Y: 2}, &Point{X: 1, Y: 2, Obj: blank}, &Point{X: 2, Y: 2}}}

	testData := ".#.\n#.#\n.#.\n"
	r := strings.NewReader(testData)

	result, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	if len(result) != len(expected) {
		t.Errorf("incorrect grid size %v; expected %v", len(result), len(expected))
		return
	}

	for x := 0; x < len(result); x++ {
		if len(result[x]) != len(expected[x]) {
			t.Errorf("incorrect grid size %v; expected %v", len(result[x]), len(expected[x]))
			return
		}

		for y := 0; y < len(result[x]); y++ {
			if result[x][y].X != expected[x][y].X {
				t.Errorf("incorrect point x value %v; expected %v", result[x][y].X, expected[x][y].X)
			}

			if result[x][y].Y != expected[x][y].Y {
				t.Errorf("incorrect point y value %v; expected %v", result[x][y].Y, expected[x][y].Y)
			}

			if result[x][y].Obj != nil && expected[x][y].Obj == nil {
				t.Errorf("incorrect obj value at (%v, %v); expected nil", x, y)
			}

			if result[x][y].Obj == nil && expected[x][y].Obj != nil {
				t.Errorf("incorrect obj value at (%v, %v); expected not nil", x, y)
			}
		}
	}
}

func TestScanViews_Small(t *testing.T) {
	testData := `.#..#
.....
#####
....#
...##
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	point := grid[4][3]
	views := ScanViews(grid, point)

	if views != 8 {
		t.Errorf("incorrect views %v; expected %v", views, 8)
	}
}

func TestScanViews_Medium(t *testing.T) {
	testData := `......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	point := grid[8][5]
	views := ScanViews(grid, point)

	if views != 33 {
		t.Errorf("incorrect views %v; expected %v", views, 33)
	}
}

func TestScanViews_Medium2(t *testing.T) {
	testData := `#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	point := grid[2][1]
	views := ScanViews(grid, point)

	if views != 35 {
		t.Errorf("incorrect views %v; expected %v", views, 35)
	}
}

func TestScanViews_Medium3(t *testing.T) {
	testData := `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	point := grid[3][6]
	views := ScanViews(grid, point)

	if views != 41 {
		t.Errorf("incorrect views %v; expected %v", views, 41)
	}
}

func TestScanViews_Large(t *testing.T) {
	testData := `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	point := grid[13][11]
	views := ScanViews(grid, point)

	if views != 210 {
		t.Errorf("incorrect views %v; expected %v", views, 210)
	}

	r2 := strings.NewReader(testData)
	grid2, err := MakeGrid(r2)
	vantage := BestView(grid2)

	if vantage == nil {
		t.Error("unexpected nil value")
		return
	}

	if vantage.X != 11 && vantage.Y != 13 {
		t.Errorf("incorrect optimal point (%v, %v); expected (%v, %v)", vantage.X, vantage.Y, 11, 13)
	}
}

func TestScanViews_FullScan(t *testing.T) {
	testData := `###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
###################
`
	r := strings.NewReader(testData)

	grid, err := MakeGrid(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err.Error())
		return
	}

	//point := grid[9][9]
	point := grid[0][18]
	views := ScanViews(grid, point)

	if views != 205 {
		t.Errorf("incorrect views %v p(%v, %v); expected %v", views, 0, 18, 205)
	}

	point = grid[18][18]
	views = ScanViews(grid, point)

	if views != 205 {
		t.Errorf("incorrect views %v p(%v, %v); expected %v", views, 18, 18, 205)
	}

	point = grid[18][0]
	views = ScanViews(grid, point)

	if views != 205 {
		t.Errorf("incorrect views %v p(%v, %v); expected %v", views, 18, 0, 205)
	}

	point = grid[0][0]
	views = ScanViews(grid, point)

	if views != 205 {
		t.Errorf("incorrect views %v p(%v, %v); expected %v", views, 0, 0, 205)
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Color ...
type Color int

const (
	// Black ...
	Black Color = 0
	// White ...
	White Color = 1
	// Transparent ...
	Transparent Color = 2
)

// Image ...
type Image struct {
	layers []ImageLayer
	width  int
	height int
}

// ImageLayer ...
type ImageLayer struct {
	data   [][]int
	digits map[int]int
}

// Render ...
func (img *Image) Render() [][]Color {
	rnd := make([][]Color, img.height)
	for u := 0; u < img.height; u++ {
		rnd[u] = make([]Color, img.width)
	}

	for i := 0; i < img.width; i++ {
		for j := 0; j < img.height; j++ {
			for _, layer := range img.layers {
				val := layer.data[j][i]
				clr := ToColor(val)
				if clr != Transparent {
					rnd[j][i] = clr
					break
				}
			}
		}
	}

	return rnd
}

// ToColor ...
func ToColor(v int) Color {
	switch v {
	case 0:
		return Black
	case 1:
		return White
	case 2:
	default:
		return Transparent
	}

	return Transparent
}

// ReadImageData ...
func ReadImageData(rd io.Reader, width int, height int) (*Image, error) {
	var err error
	r := bufio.NewReader(rd)

	img := Image{layers: []ImageLayer{}, width: width, height: height}

	layer := ImageLayer{data: [][]int{}, digits: make(map[int]int)}
	row := []int{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if len(row) > 0 {
				layer.data = append(layer.data, row)
			}

			if len(layer.data) > 0 {
				img.layers = append(img.layers, layer)
			}
			break
		}

		i, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, err
		}

		cnt, ok := layer.digits[i]
		if !ok {
			layer.digits[i] = 1
		} else {
			layer.digits[i] = cnt + 1
		}

		row = append(row, i)
		if len(row) == width {
			layer.data = append(layer.data, row)
			row = []int{}
		}

		if len(layer.data) == height {
			img.layers = append(img.layers, layer)
			layer = ImageLayer{data: [][]int{}, digits: make(map[int]int)}
		}
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	return &img, nil
}

func main() {
	var data string
	var img *Image

	file, err := os.Open("./data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	img, err = ReadImageData(strings.NewReader(data), 25, 6)
	if err != nil {
		log.Fatal(err)
	}

	least := int(^uint(0) >> 1)
	li := 0
	for i, layer := range img.layers {
		val, ok := layer.digits[0]
		if ok && val < least {
			least = val
			li = i
		}
	}

	o, _ := img.layers[li].digits[1]
	t, _ := img.layers[li].digits[2]

	final := o * t

	println(fmt.Sprintf("final value: %v", final))

	rendering := img.Render()
	for i := 0; i < len(rendering); i++ {
		r := rendering[i]
		line := ""
		for j := 0; j < len(r); j++ {
			switch r[j] {
			case White:
				line += "@"
			case Transparent:
				line += "."
			case Black:
				line += " "
			}
		}
		println(line)
	}
}

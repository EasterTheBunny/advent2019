package main

import (
	"strings"
	"testing"
)

func TestReadImageData(t *testing.T) {
	data := "123456789012"
	buf := strings.NewReader(data)
	img, err := ReadImageData(buf, 3, 2)
	if err != nil {
		t.Errorf("read error occurred: %s", err.Error())
	}

	if len(img.layers) != 2 {
		t.Errorf("incorrect number of layers %v; expected %v", len(img.layers), 2)
		return
	}
}

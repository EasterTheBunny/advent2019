package main

import (
	"strings"
	"testing"
)

func TestInSequence(t *testing.T) {
	inputB := []string{"A", "Z", "1", "4"}
	inputA := []string{"B", "Y", "2", "3"}
	expected := []bool{false, true, false, true}

	for i, b := range inputB {
		value := InSequence(b, inputA[i])

		if value != expected[i] {
			t.Errorf("incorrect sequence result %v; expected %v", value, expected[i])
		}
	}
}

func TestNextGroup(t *testing.T) {
	inputs := []string{"445", "333", "22", "213", "5556", "111122"}
	inputsS := []int{0, 0, 0, 1, 0, 4}
	expectedI := []int{2, -1, -1, 2, 3, -1}
	expectedL := []int{2, 3, 2, 1, 3, 2}

	for i, input := range inputs {
		seq, next := NextGroup(strings.Split(input, ""), inputsS[i])

		if next != expectedI[i] {
			t.Errorf("incorrect next index result %v for input %v; expected %v", seq, input, expectedI[i])
		}

		if len(seq) != expectedL[i] {
			t.Errorf("incorrect length result %v for input %v; expected %v", len(seq), input, expectedL[i])
		}
	}
}

func TestValidPassword(t *testing.T) {
	inputs := []string{"111111", "223450", "123789", "112233", "123444", "111122"}
	expected := []bool{false, false, false, true, false, true}

	for i, input := range inputs {
		value := ValidPassword(input)

		if value != expected[i] {
			t.Errorf("incorrect valid result %v for input %v; expected %v", value, input, expected[i])
		}
	}
}

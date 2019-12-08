package main

import (
	"testing"
)

func TestRunSetting(t *testing.T) {
	codes := [][]int{
		[]int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
		[]int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
		[]int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0}}
	setting := [][]int{
		[]int{4, 3, 2, 1, 0},
		[]int{0, 1, 2, 3, 4},
		[]int{1, 0, 4, 3, 2}}

	expected := []int{43210, 54321, 65210}

	for i, set := range codes {
		result := RunSetting(setting[i], set)

		if result != expected[i] {
			t.Errorf("incorrect signal %v; expected %v", result, expected[i])
		}
	}
}

func TestSettingsList(t *testing.T) {
	expected := 3125

	result := SettingsList(0, 4)

	if len(result) != expected {
		t.Errorf("incorrect possible settings %v; expected %v", len(result), expected)
	}
}

func TestValidSetting(t *testing.T) {
	inputs := [][]int{
		[]int{0, 1, 2, 3, 4},
		[]int{2, 2, 3, 4, 1}}

	expected := []bool{true, false}

	for i, input := range inputs {
		valid := ValidSetting(input)

		if valid != expected[i] {
			t.Errorf("incorrect validity test %v on index %v; expected %v", valid, i, expected[i])
		}
	}
}

func TestMaxSetting(t *testing.T) {
	codes := [][]int{
		[]int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
		[]int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
		[]int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0}}

	expected := []int{43210, 54321, 65210}

	for i, set := range codes {
		max, _ := MaxSetting(set)

		if max != expected[i] {
			t.Errorf("incorrect signal %v; expected %v", max, expected[i])
		}
	}
}

package main

import (
	"testing"
)

func TestRunSetting_SimpleMode(t *testing.T) {
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

func TestMaxSetting_SimpleMode(t *testing.T) {
	codes := [][]int{
		[]int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0},
		[]int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0},
		[]int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0}}

	expected := []int{43210, 54321, 65210}

	for i, set := range codes {
		max, _ := MaxSetting(SimpleMode, set)

		if max != expected[i] {
			t.Errorf("incorrect signal %v; expected %v", max, expected[i])
		}
	}
}

func TestRunSetting_FeedbackMode(t *testing.T) {
	codes := [][]int{
		[]int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
		[]int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}}
	setting := [][]int{
		[]int{9, 8, 7, 6, 5},
		[]int{9, 7, 8, 5, 6}}

	expected := []int{139629729, 18216}

	for i, set := range codes {
		result := RunSetting(setting[i], set)

		if result != expected[i] {
			t.Errorf("incorrect signal %v; expected %v", result, expected[i])
		}
	}
}

func TestMaxSetting_FeedbackMode(t *testing.T) {
	codes := [][]int{
		[]int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
		[]int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}}

	expected := []int{139629729, 18216}

	for i, set := range codes {
		max, _ := MaxSetting(FeedbackMode, set)

		if max != expected[i] {
			t.Errorf("incorrect signal %v; expected %v", max, expected[i])
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

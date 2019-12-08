package main

import (
	"2019/internal/intcode"
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// RunSetting ...
func RunSetting(setting []int, set []int) int {
	amps := make([][]int, len(setting))
	for i := 0; i < len(amps); i++ {
		a := make([]int, len(set))
		copy(a, set)
		amps[i] = a
	}

	output := 0
	for i, s := range amps {
		sl := strconv.Itoa(setting[i]) + "\n" + strconv.Itoa(output) + "\n"
		inBytes := []byte(sl)

		br := bytes.NewReader(inBytes)
		reader := bufio.NewReader(br)
		writer := bytes.NewBuffer([]byte{})

		intcode.Process(reader, writer, 0, s)

		r := strings.Split(strings.TrimRight(fmt.Sprintf("%v", writer), "\n"), "\n")
		j, err := strconv.Atoi(r[len(r)-1])
		if err != nil {
			panic("non-int output")
		}
		output = j
	}

	return output
}

// ValidSetting ...
func ValidSetting(setting []int) bool {
	codes := make(map[int]bool)

	for _, i := range setting {
		_, ok := codes[i]
		if ok {
			return false
		}
		codes[i] = true
	}
	return true
}

// SettingsList ...
func SettingsList(min int, max int) [][]int {
	list := [][]int{}

	for i := min; i <= max; i++ {
		for j := min; j <= max; j++ {
			for k := min; k <= max; k++ {
				for l := min; l <= max; l++ {
					for m := min; m <= max; m++ {
						item := []int{i, j, k, l, m}
						list = append(list, item)
					}
				}
			}
		}
	}

	return list
}

// MaxSetting ...
func MaxSetting(codes []int) (int, []int) {
	max := 0
	s := []int{0, 0, 0, 0, 0}
	allSettings := SettingsList(0, 4)

	for _, setting := range allSettings {
		if ValidSetting(setting) {
			result := RunSetting(setting, codes)
			if result > max {
				s = setting
				max = result
			}
		}
	}

	return max, s
}

func main() {
	codes := intcode.ReadCodes("./commands.txt")

	max, _ := MaxSetting(codes)

	println(fmt.Sprintf("max value: %v", max))
	// wrong answer: 76332
	// wrong answer: 25076
}

package main

import (
	"2019/internal/intcode"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RunMode ...
type RunMode int

const (
	// FeedbackMode ...
	FeedbackMode RunMode = 1
	// SimpleMode ...
	SimpleMode RunMode = 2
)

// Amplifier ...
type Amplifier struct {
	Data   []int
	Input  io.Reader
	Output io.Writer
}

// RunSetting ...
func RunSetting(setting []int, set []int) int {
	amps := make([]*Amplifier, len(setting))

	startReader, endWriter := io.Pipe()
	for i := 0; i < len(amps); i++ {
		a := make([]int, len(set))
		copy(a, set)
		amp := Amplifier{
			Data: a}

		if i == 0 {
			amp.Input = startReader
		} else {
			reader, writer := io.Pipe()
			amps[i-1].Output = writer
			amp.Input = reader
		}
		amps[i] = &amp
	}
	amps[len(amps)-1].Output = endWriter

	done := make(chan bool, len(setting)-1)

	for _, s := range amps {
		go func(amp *Amplifier) {
			intcode.Process(amp.Input, amp.Output, 0, amp.Data)
			done <- true
		}(s)
	}

	for i, val := range append(setting[1:], setting[0]) {
		fmt.Fprint(amps[i].Output, fmt.Sprintf("%v\n", val))
	}

	fmt.Fprint(amps[len(amps)-1].Output, "0\n")

	for i := 0; i < len(setting)-1; i++ {
		<-done
	}

	reader := bufio.NewReader(startReader)
	text, _ := reader.ReadString('\n')

	j, err := strconv.Atoi(strings.TrimRight(text, "\n"))
	if err != nil {
		panic("non-int output")
	}

	return j
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
func MaxSetting(mode RunMode, codes []int) (int, []int) {
	max := 0
	s := []int{0, 0, 0, 0, 0}
	allSettings := [][]int{}

	switch mode {
	case SimpleMode:
		allSettings = SettingsList(0, 4)
	case FeedbackMode:
		allSettings = SettingsList(5, 9)
	}

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

	max, _ := MaxSetting(SimpleMode, codes)
	println(fmt.Sprintf("max value simple mode: %v", max))
	max, _ = MaxSetting(FeedbackMode, codes)
	println(fmt.Sprintf("max value feedback mode: %v", max))
}

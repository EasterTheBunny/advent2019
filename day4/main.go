package main

import (
	"fmt"
	"strconv"
	"strings"
)

// InSequence ...
func InSequence(b string, a string) bool {
	return b > a
}

// NextGroup ...
func NextGroup(seq []string, start int) ([]string, int) {
	if len(seq) == 0 || start+1 == len(seq) {
		return []string{}, -1
	}

	a := seq[start]
	r := []string{a}
	for i := start + 1; i < len(seq); i++ {
		b := seq[i]
		if a != b {
			return r, i
		}
		r = append(r, b)
	}

	return r, -1
}

// ValidPassword validates a password by the following rules:
// must contain two matching adjacent numbers
// numbers must increase or stay the same toward the end
// must be a 6 digit number
func ValidPassword(p string) bool {
	chars := strings.Split(p, "")

	if len(chars) != 6 {
		return false
	}

	validGrouping := false
	validSequence := true

	i := 0
	for sub, i := NextGroup(chars, i); validSequence; sub, i = NextGroup(chars, i) {
		if len(sub) == 2 {
			validGrouping = true
		}

		if i > 0 && !InSequence(chars[i], sub[len(sub)-1]) {
			validSequence = false
		}

		if i < 0 {
			break
		}
	}

	return validGrouping && validSequence
}

func main() {
	lowerLimit := 256310
	upperLimit := 732736
	valids := []string{}

	for i := lowerLimit; i <= upperLimit; i++ {
		s := strconv.Itoa(i)
		v := ValidPassword(s)
		if v {
			valids = append(valids, s)
		}
	}
	println(fmt.Sprintf("valid passwords: %v", len(valids)))
}

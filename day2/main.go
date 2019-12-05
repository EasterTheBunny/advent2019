package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func add(i, j int) int {
	return i + j
}

func mult(i, j int) int {
	return i * j
}

func op(position int, v []int) (code, i, j, rPos int) {
	code = v[position]
	if code == 99 {
		i = 0
		j = 0
		rPos = 0
		return
	}

	i = v[position+1]
	j = v[position+2]
	rPos = v[position+3]
	return
}

func switchOps(position int, v []int) []int {
	code, i, j, rPos := op(position, v)

	switch code {
	case 1:
		v[rPos] = add(v[i], v[j])
		return switchOps(position+4, v)
	case 2:
		v[rPos] = mult(v[i], v[j])
		return switchOps(position+4, v)
	case 99:
		return v
	default:
		return v
	}
}

func main() {
	base := []int{}

	file, err := os.Open("./opcodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ",")
		for _, v := range values {
			vi, _ := strconv.Atoi(v)
			base = append(base, vi)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	a := make([]int, len(base))
	copy(a, base)
	// replace position 1 with the value 12
	a[1] = 12
	// replace position 2 with the value 2
	a[2] = 2

	a = switchOps(0, a)

	println(fmt.Sprintf("basic result: %v", a[0]))

	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			b := make([]int, len(base))
			copy(b, base)
			b[1] = i
			b[2] = j
			b = switchOps(0, b)
			if b[0] == 19690720 {
				println(fmt.Sprintf("noun: %v; verb: %v", i, j))
				println(fmt.Sprintf("result: %v", (100*i)+j))
				break
			}
		}
	}
}

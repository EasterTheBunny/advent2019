package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var a = []int{}

func p(m int) int {
	return m/3 - 2
}

func q(m int) int {
	t := p(m)
	if t < 0 {
		return 0
	}
	return t + q(t)
}

func main() {
	file, err := os.Open("./modules.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v, _ := strconv.Atoi(scanner.Text())
		a = append(a, v)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	d := 0
	e := 0
	for _, value := range a {
		d += p(value)
		e += q(value)
	}
	println(fmt.Sprintf("sum: %v", d))
	println(fmt.Sprintf("recursive sum: %v", e))
}

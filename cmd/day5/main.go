package main

import (
	"2019/internal/intcode"
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	intcode.Process(reader, os.Stdout, 0, intcode.ReadCodes("./instructions.txt"))
}

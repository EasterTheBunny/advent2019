package main

import (
	"2019/internal/intcode"
	"os"
)

func main() {
	codes := intcode.ReadCodes("./codes.txt")

	intcode.Process(os.Stdin, os.Stdout, 0, codes)
}

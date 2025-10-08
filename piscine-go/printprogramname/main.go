package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arguments := os.Args[0]
	aR := []rune(arguments)

	for i, v := range aR {
		if i > 1 {
			z01.PrintRune(rune(v))
		}
	}
	z01.PrintRune('\n')
}

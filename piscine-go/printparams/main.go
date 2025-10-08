package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arguments := os.Args
	countArguments := 0

	for i := range arguments {
		countArguments = i
	}
	for i := 1; i <= countArguments; i++ {
		for _, v := range arguments[i] {
			z01.PrintRune(v)
		}
		z01.PrintRune('\n')
	}
}

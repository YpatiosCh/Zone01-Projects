package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	arguments := os.Args
	if len(arguments) <= 1 {
		return // No arguments provided, exit
	}

	table := arguments[1:]

	if table[0] == "--upper" {
		tableUpp := table[1:]
		for _, numStr := range tableUpp {
			num := basicAtoi(numStr)
			if num >= 1 && num <= 26 {
				z01.PrintRune(rune(num + 64))
			} else {
				z01.PrintRune(' ')
			}
		}
	} else {
		for _, numStr := range table {
			num := basicAtoi(numStr)
			if num >= 1 && num <= 26 {
				z01.PrintRune(rune(num + 96))
			} else {
				z01.PrintRune(' ')
			}
		}
	}
	z01.PrintRune('\n')
}

func basicAtoi(s string) int {
	var num int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		} else {
			return 0 // Invalid character, return 0
		}
	}
	return num
}

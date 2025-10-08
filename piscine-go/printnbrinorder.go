package piscine

import "github.com/01-edu/z01"

func PrintNbrInOrder(n int) {
	if n == 0 {
		z01.PrintRune('0')
		return
	}

	// Array to store the count of each digit
	digitCount := [10]int{}

	// Count the occurrence of each digit in the number
	for n > 0 {
		digit := n % 10
		digitCount[digit]++
		n /= 10
	}

	// Print the digits in ascending order
	for i := 0; i < 10; i++ {
		for j := 0; j < digitCount[i]; j++ {
			// Use z01.PrintRune to print each digit individually
			z01.PrintRune(rune('0' + i))
		}
	}
}

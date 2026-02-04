package main

import (
	"github.com/01-edu/z01"
)

// SetPoint sets the values of two integers to 42 and 21 respectively.
// It takes pointers to two integers as arguments and updates the values they point to.
func SetPoint(x, y *int) {
	*x = 42 // Set the value pointed by x to 42
	*y = 21 // Set the value pointed by y to 21
}

func main() {
	var x, y int // Declare two integer variables: x and y

	SetPoint(&x, &y) // Call the SetPoint function passing the addresses of x and y

	// Print 'x = ' followed by the value of x, then a comma and a space,
	// followed by 'y = ' and the value of y, and finally a new line character using the z01 package.
	printString("x = ")
	printInt(x)
	printString(", y = ")
	printInt(y)
	printRune('\n')
}

// printString prints a string using the z01 package.
func printString(s string) {
	for _, char := range s {
		z01.PrintRune(char)
	}
}

// printRune prints a rune using the z01 package.
func printRune(r rune) {
	z01.PrintRune(r)
}

// printInt prints an integer using the z01 package.
func printInt(n int) {
	if n == 0 {
		z01.PrintRune('0')
		return
	}
	if n < 0 {
		z01.PrintRune('-')
		n = -n
	}
	digits := []rune{}
	for n > 0 {
		d := n % 10
		digits = append(digits, rune(d+48)) // Convert int to rune and add 48 to get the ASCII representation of the digit.
		n /= 10
	}
	// Print digits in reverse order
	for i := len(digits) - 1; i >= 0; i-- {
		z01.PrintRune(digits[i])
	}
}

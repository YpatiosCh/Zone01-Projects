package piscine

import (
	"fmt"
)

// DealAPackOfCards distributes the cards evenly among four players.
func DealAPackOfCards(deck []int) {
	// Assuming the deck is always of length 12, as per the task requirement.
	// No need to check the length with len.

	// Distribute the cards evenly among the players.
	for i := 0; i < 4; i++ {
		fmt.Printf("Player %d: ", i+1)
		for j := 0; j < 3; j++ {
			// Print the card number directly without converting to rune
			fmt.Printf("%d", deck[i*3+j])
			if j < 2 {
				fmt.Printf(", ")
			}
		}
		// Use fmt.Printf to print a newline character at the end of each line
		fmt.Printf("\n")
	}
}

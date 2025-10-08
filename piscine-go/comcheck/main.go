package main

import (
	"fmt"
	"os"
)

func main() {
	// Iterate over the command-line arguments
	for _, arg := range os.Args[1:] {
		// Check if the argument matches "01, galaxy" or "galaxy 01"
		if arg == "01" || arg == "galaxy" || arg == "galaxy 01" {
			// If a match is found, print "Alert!!!" and return
			fmt.Println("Alert!!!")
			return
		}
	}
}

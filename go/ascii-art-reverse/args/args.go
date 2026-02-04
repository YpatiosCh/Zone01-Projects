package args

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// PrintReverseUsage prints to terminal how to use the program.
func PrintReverseUsage() {
	fmt.Println("\nUsage: go run . [OPTION]")
	fmt.Println("\nEX: go run . --reverse=<fileName>")
}

// HandleArgs handles the arguments as needed and returns the flag input as a pointer to a string.
// If there is an error, it returns the error.
func HandleArgs() (*string, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("\nno arguments provided")
	}

	revflag := os.Args[1:]

	err := ValidateReverseFlag(revflag)
	if err != nil {
		return nil, err
	}

	reverseFlag := flag.String("reverse", "", "Specify file you want to reverse")
	flag.Parse()

	return reverseFlag, nil
}

// ValidateReverseFlag will validate if the --reverse flag is being used correctly.
func ValidateReverseFlag(t []string) error {
	if len(t) > 1 {
		return fmt.Errorf("\ninvalid syntax for reverse flag. please use '--reverse=<filename>' format")
	}

	if strings.HasPrefix(t[0], "--reverse") && !strings.Contains(t[0], "=") {
		return fmt.Errorf("\ninvalid syntax for reverse flag. please use '--reverse=<filename>' format")
	}
	if strings.HasPrefix(t[0], "--reverse=") && !strings.HasSuffix(t[0], ".txt") {
		return fmt.Errorf("\ninvalid file. can only reverse txt files")

	}
	if !strings.HasPrefix(t[0], "--reverse=") {
		return fmt.Errorf("\ninvalid argument. please use '--reverse=<filename>'")
	}
	return nil
}

package ascii

import (
	"errors"
	"flag"
	"fmt"
)

// HandleArgs processes and validates command-line arguments, returning required parameters or an error if invalid
func HandleArgs() (string, string, string, string, string, string, error) {
	// Validate the --color= flag format
	if err := ValidateColorFlag(); err != nil {
		return "", "", "", "", "", "", err
	}

	// Validate the --align= flag format
	if err := ValidateAlignFlag(); err != nil {
		return "", "", "", "", "", "", err
	}

	// Validate the --output= flag format
	if err := ValidateOutputFlag(); err != nil {
		return "", "", "", "", "", "", err
	}

	// Define and parse color and output flag
	colorFlag := flag.String("color", "white", "Specify the output color (e.g., red, green, yellow)")
	alignFlag := flag.String("align", "left", "Specify alignment: center, left, right, justify")
	outputFlag := flag.String("output", "", "Specify an output file (e.g., banner.txt)")
	flag.Parse()

	err := ValidateAlignment(*alignFlag)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	err = HandleFlags(*colorFlag, *alignFlag, *outputFlag)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	// Check if we have enough arguments
	if flag.NArg() < 1 {
		return "", "", "", "", "", "", errors.New("\nnot enough arguments")
	}

	// Default values
	substring := flag.Arg(0)
	font := "./fonts/standard.txt"
	var fullInput string

	switch flag.NArg() {
	case 1:
		// Single argument, treat as both substring and input
		fullInput = flag.Arg(0)
	case 2:
		if IsSubstring(flag.Arg(1), flag.Arg(0)) {
			// Two arguments, with second as full input
			substring = flag.Arg(0)
			fullInput = flag.Arg(1)
		} else if !IsSubstring(flag.Arg(1), flag.Arg(0)) && !IsFont(flag.Arg(1)) { // if first arg is not a substring and the second arg is not a fontName
			substring = flag.Arg(0)
			fullInput = flag.Arg(1)
		} else {
			// Second argument specifies a font
			fullInput = flag.Arg(0)
			font = "./fonts/" + flag.Arg(1) + ".txt"
		}

	case 3:
		// Three arguments: substring, input, and font name
		substring = flag.Arg(0)
		fullInput = flag.Arg(1)
		font = "./fonts/" + flag.Arg(2) + ".txt"
	default:
		PrintUsage()
		return "", "", "", "", "", "", errors.New("\ninvalid number of arguments")
	}

	return substring, fullInput, *colorFlag, font, *outputFlag, *alignFlag, nil
}

func HandleFlags(colorFlag, alignFlag, outputFlag string) error {
	if outputFlag != "" && colorFlag != "white" {
		return fmt.Errorf("\ncannot use color in a text file")
	} else if outputFlag != "" && alignFlag != "left" {
		return fmt.Errorf("\ncannot use alignment in a text file")
	}
	return nil
}

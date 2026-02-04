package main

import (
	"ascii-art-output/tools"
	"os"
)

func main() {
	var input string
	var outputFile string
	font := "./fonts/standard.txt"

	argsLen := len(os.Args)

	// Check for the number of arguments
	if argsLen < 2 {
		tools.PrintUsage()
		return
	} else if argsLen == 2 {
		// Check if the output flag is provided
		isOutputFlag := tools.ToFileOrNotToFile(os.Args[1])

		if isOutputFlag {
			tools.PrintUsage()
			return
		}
		// No errors, proceed with printing to terminal
		input = os.Args[1]
		finalOutput := tools.CreateOutput(input, font)

		tools.WriteToTerminal(finalOutput)

	} else if argsLen == 3 {
		// Check if we are writing to a file or printing to terminal
		isFile := tools.ToFileOrNotToFile(os.Args[1])

		if !isFile {
			// Print to terminal with specified font
			input = os.Args[1]
			font = "./fonts/" + os.Args[2] + ".txt"
			finalOutput := tools.CreateOutput(input, font)
			tools.WriteToTerminal(finalOutput)
		} else {
			// Write to file with standard font
			input = os.Args[2]
			finalOutput := tools.CreateOutput(input, font)
			outputFile = tools.ExtractFileName(os.Args[1])
			tools.WriteToFile(finalOutput, outputFile)
		}
	} else if argsLen == 4 {
		// Handle case with both output file and specified font
		isFile := tools.ToFileOrNotToFile(os.Args[1])

		if !isFile {
			tools.PrintUsage()
			return
		}
		// Write to file with specified font
		outputFile = tools.ExtractFileName(os.Args[1])
		input = os.Args[2]
		font = "./fonts/" + os.Args[3] + ".txt"
		finalOutput := tools.CreateOutput(input, font)
		tools.WriteToFile(finalOutput, outputFile)
	} else {
		// If arguments exceed the expected count, print usage
		tools.PrintUsage()
		return
	}
}

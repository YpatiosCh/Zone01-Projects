package ascii

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ValidateOutputFlag ensures the output flag has a valid format (ends with .txt and uses `=` syntax).
func ValidateOutputFlag() error {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--output") && !strings.Contains(arg, "=") && !strings.HasSuffix(arg, ".txt") {
			return fmt.Errorf("\ninvalid syntax for --color flag. please use '--color=colorname' format")
		}
	}
	return nil
}

// WriteToFile writes ASCII art to the specified output file
// Uses stripAnsi helper function to strip the color of the output so it can print to file
func WriteToFile(outputFile string, lines []string) {
	// Create or overwrite the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		PrintUsage()
	}
	defer outFile.Close()

	for _, line := range lines {
		line = stripAnsi(line)
		outFile.WriteString(line + "\n")
	}

}

// Helper function to remove ANSI color codes from strings
func stripAnsi(input string) string {
	// Strip any ANSI escape sequences
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

package ascii

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

// determineColor gets the color needed from the flag --color=
func determineColor(colorFlag *string) string {
	// Determine ANSI color code based on the specified color
	var color string
	switch *colorFlag {
	case "red":
		color = red
	case "green":
		color = green
	case "yellow":
		color = yellow
	case "blue":
		color = blue
	case "magenta":
		color = magenta
	case "cyan":
		color = cyan
	default:
		color = white // Default to white
	}
	return color
}

// Function to determine which characters should be colored
func determineSubToColor(input string, substring string) []bool {
	colorFlags := make([]bool, len(input))
	// Find occurrences of the substring
	for i := 0; i <= len(input)-len(substring); i++ {
		if input[i:i+len(substring)] == substring {
			for j := 0; j < len(substring); j++ {
				colorFlags[i+j] = true
			}
		}
	}
	return colorFlags
}

// Validate the color flag to ensure it uses "=" syntax
func ValidateColorFlag() error {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--color") && !strings.Contains(arg, "=") {
			return fmt.Errorf("\ninvalid syntax for --color flag. please use '--color=colorname' format")
		}
	}
	return nil
}

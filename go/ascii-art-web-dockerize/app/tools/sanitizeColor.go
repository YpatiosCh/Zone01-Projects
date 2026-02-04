package tools

import "regexp"

// Regular expression to match valid hex color values
var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// sanitizeColor ensures that the color is in the form of a valid hex code
func SanitizeColor(color string) string {
	if hexColorPattern.MatchString(color) {
		return color // Valid color
	}
	return "#000000" // Default to black if invalid
}
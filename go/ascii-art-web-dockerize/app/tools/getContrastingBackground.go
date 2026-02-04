package tools

import "fmt"

// getContrastingBackground determines a contrasting background color
func GetContrastingBackground(color string) string {
    // Convert hex color to RGB components
    var r, g, b int
    fmt.Sscanf(color, "#%02x%02x%02x", &r, &g, &b)

    // Calculate the brightness using a luminance formula
    brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

    // Return white for dark colors and black for light colors
    if brightness < 128 {
        return "#FFFFFF" // Light background for dark text
    }
    return "#000000" // Dark background for light text
}
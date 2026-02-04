package utils

import (
	"groupie-tracker/models"
	"strings"
)

// formatLocations processes locations to make them more readable
func FormatLocations(data *models.LocationsData) {
	for i, loc := range data.Locations {
		data.Locations[i] = FormatLocation(loc)
	}
}

// formatRelations processes both locations and dates in relations data
func FormatRelations(data *models.RelationsData) {
	formatted := make(map[string][]string)

	// Format each location key and its dates
	for loc, dates := range data.DatesLocations {
		formattedLoc := FormatLocation(loc)
		formatted[formattedLoc] = dates
	}

	data.DatesLocations = formatted
}

// formatLocation handles the formatting of a single location string
func FormatLocation(loc string) string {
	// Replace underscores and hyphens with spaces and commas
	loc = strings.ReplaceAll(loc, "-", ", ")
	loc = strings.ReplaceAll(loc, "_", " ")

	// Split into words for capitalization
	words := strings.Split(loc, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
		case "usa":
			words[i] = strings.ToUpper(words[i])
		case "uk":
			words[i] = strings.ToUpper(words[i])
		default:
			words[i] = strings.ToTitle(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

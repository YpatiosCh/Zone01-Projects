package filter

import (
	"groupie-tracker/models"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Filters represents the criteria for filtering artists
type Filters struct {
	CreationDateMin int      // Minimum creation date of the artist (Year)
	CreationDateMax int      // Maximum creation date of the artist (Year)
	FirstAlbumMin   int      // Minimum year for the artist's first album
	FirstAlbumMax   int      // Maximum year for the artist's first album
	MemberMin       int      // Minimum number of members in the artist's group
	MemberMax       int      // Maximum number of members in the artist's group
	Locations       []string // List of locations to filter by
}

// LocationsByCountry represents a structured location with city and country
type LocationsByCountry struct {
	Country string   // Country name
	Cities  []string // List of cities within the country
}

// InRange checks if a value is within a given range [min, max]
func InRange(value, min, max int) bool {
	return value >= min && value <= max
}

// GetIntParam retrieves an integer parameter from the HTTP request, returning a default value if not found
func GetIntParam(r *http.Request, param string, defaultVal int) int {
	val := r.FormValue(param)
	if val == "" {
		return defaultVal
	}

	// Convert the parameter to an integer, and return default value if conversion fails
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}

// InAlbumRange checks if the album's year is within the specified range (inclusive)
func InAlbumRange(albumDate string, min, max int) bool {
	// Extract the year part from the album's date (format: "DD-MM-YYYY")
	parts := strings.Split(albumDate, "-")
	if len(parts) != 3 {
		return false
	}

	// Convert the year (last part of the date) to an integer and check if it's within the range
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return false
	}

	return year >= min && year <= max
}

// CheckLocations verifies if an artist's locations match all the selected filter locations
func CheckLocations(artistLocations []string, filterLocations []string) bool {
	// If no locations are specified in the filter, return true (no filtering applied)
	if len(filterLocations) == 0 {
		return true
	}

	// Check if each filter location exists in the artist's locations
	for _, filterLoc := range filterLocations {
		locationFound := false
		for _, artistLoc := range artistLocations {
			// Split artist location into city and country, and check if the city matches
			parts := strings.Split(artistLoc, ",")
			city := strings.TrimSpace(parts[0])
			if city == filterLoc {
				locationFound = true
				break
			}
		}
		// If any filter location is not found in the artist's locations, return false
		if !locationFound {
			return false
		}
	}
	return true
}

// GetUniqueLocations retrieves a sorted list of unique locations (cities and countries) from a list of artists
func GetUniqueLocations(artists []models.Artist) []LocationsByCountry {
	// Map to store cities by country
	locationMap := make(map[string]map[string]bool)

	// Loop through all artists to collect their locations
	for _, artist := range artists {
		for _, location := range artist.LocationsData.Locations {
			// Split the location into city and country
			parts := strings.Split(location, ",")
			if len(parts) == 2 {
				city := strings.TrimSpace(parts[0])
				country := strings.TrimSpace(parts[1])

				// Initialize country map if it doesn't exist
				if _, exists := locationMap[country]; !exists {
					locationMap[country] = make(map[string]bool)
				}
				// Add the city to the country's map
				locationMap[country][city] = true
			}
		}
	}

	// Convert the map to a sorted slice of LocationsByCountry
	var locations []LocationsByCountry
	for country, cities := range locationMap {
		cityList := make([]string, 0, len(cities))
		for city := range cities {
			cityList = append(cityList, city)
		}
		// Sort cities alphabetically
		sort.Strings(cityList)
		locations = append(locations, LocationsByCountry{
			Country: country,
			Cities:  cityList,
		})
	}

	// Sort the list of locations by country name
	sort.Slice(locations, func(i, j int) bool {
		return locations[i].Country < locations[j].Country
	})

	return locations
}

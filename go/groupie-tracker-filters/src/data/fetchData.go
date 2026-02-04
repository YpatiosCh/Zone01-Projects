package data

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/utils"
	"io"
	"net/http"
	"sync"
)

const defaultArtistAPIEndpoint = "https://groupietrackers.herokuapp.com/api/artists"

// FetchArtists retrieves artist data from the API and their related information.
// This is the main function that orchestrates the entire data fetching process.
// It first fetches the basic artist information, then concurrently fetches related
// data (locations, concert dates, and relations) for each artist.
func FetchArtists() ([]models.Artist, error) {
	fmt.Println("Fetching artists...")
	// Get the API URL from environment variable "myData"
	resp, err := http.Get(defaultArtistAPIEndpoint)
	if err != nil {
		// If the HTTP request fails, return a wrapped error with context
		return nil, fmt.Errorf("error fetching artists: %v", err)
	}
	// Ensure the response body is closed after we're done with it
	// defer ensures this happens even if the function returns early
	defer resp.Body.Close()
	// Read the entire response body into memory
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If reading the response body fails, return a wrapped error
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	// Parse the JSON response into our Artist structs
	var artists []models.Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		// If JSON parsing fails, return a wrapped error
		return nil, fmt.Errorf("error unmarshaling artists: %v", err)
	}
	// Set up concurrent fetching of related data for each artist
	// Create a WaitGroup to track when all goroutines complete
	var wg sync.WaitGroup
	// Create a buffered error channel to collect errors from goroutines
	// Buffer size matches number of artists to prevent goroutine leaks
	errChan := make(chan error, len(artists))
	// Launch a goroutine for each artist to fetch their related data
	for i := range artists {
		// Increment WaitGroup counter before launching goroutine
		wg.Add(1)
		// Launch goroutine to fetch related data
		go func(artist *models.Artist) {
			// Ensure WaitGroup is decremented when goroutine completes
			defer wg.Done()
			// Fetch related data for this artist
			if err := fetchRelatedDataAsync(artist); err != nil {
				// If an error occurs, send it to the error channel
				errChan <- fmt.Errorf("error for artist %d: %v", artist.ID, err)
			}
		}(&artists[i]) // Pass address of artist to avoid loop variable capture issues
	}
	// Wait for all goroutines to complete
	wg.Wait()
	// Close error channel to prevent resource leaks
	close(errChan)
	// Check if any errors occurred during the concurrent fetching
	for err := range errChan {
		if err != nil {
			// Return the first error encountered
			return nil, err
		}
	}
	// Return the fully populated artist slice
	return artists, nil
}

// fetchRelatedDataAsync concurrently fetches all related data for a single artist.
// This includes locations, concert dates, and relations data.
// The function launches three goroutines simultaneously to fetch data in parallel.
func fetchRelatedDataAsync(artist *models.Artist) error {
	// Set up concurrency primitives
	// WaitGroup to track completion of all three goroutines
	var wg sync.WaitGroup
	// Buffered channel for collecting errors from goroutines
	// Buffer size is 3 to match number of goroutines
	errChan := make(chan error, 3)
	// Prepare to launch three goroutines
	wg.Add(3)
	// Launch goroutine for fetching locations data
	go func() {
		// Ensure WaitGroup is decremented when goroutine completes
		defer wg.Done()
		// Fetch and parse locations data
		if err := fetchAndUnmarshal(artist.Locations, &artist.LocationsData); err != nil {
			errChan <- fmt.Errorf("locations error: %v", err)
		}
	}()
	// Launch goroutine for fetching concert dates data
	go func() {
		defer wg.Done()
		// Fetch and parse concert dates data
		if err := fetchAndUnmarshal(artist.ConcertDates, &artist.ConcertDatesData); err != nil {
			errChan <- fmt.Errorf("dates error: %v", err)
		}
	}()
	// Launch goroutine for fetching relations data
	go func() {
		defer wg.Done()
		// Fetch and parse relations data
		if err := fetchAndUnmarshal(artist.Relations, &artist.RelationsData); err != nil {
			errChan <- fmt.Errorf("relations error: %v", err)
		}
	}()
	// Wait for all three goroutines to complete
	wg.Wait()
	// Close error channel to prevent resource leaks
	close(errChan)
	// Check if any errors occurred during the concurrent fetching
	for err := range errChan {
		if err != nil {
			// Return the first error encountered
			return err
		}
	}
	// Return nil if all operations succeeded
	return nil
}

// fetchAndUnmarshal performs an HTTP GET request and unmarshals the JSON response.
// This is a utility function used by fetchRelatedDataAsync to handle the common
// pattern of fetching JSON data and parsing it into a struct.
func fetchAndUnmarshal[T *models.DatesData | *models.LocationsData | *models.RelationsData](url string, target T) error {
	// Perform HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		// Return error if the HTTP request fails
		return err
	}
	// Ensure response body is closed after we're done
	defer resp.Body.Close()
	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Return error if reading the body fails
		return err
	}
	// Parse JSON into the target struct
	// Return any error that occurs during unmarshaling
	if err = json.Unmarshal(body, target); err != nil {
		return err
	}
	// Apply specific formatting based on type
	switch v := any(target).(type) {
	case *models.LocationsData:
		utils.FormatLocations(v)
	case *models.RelationsData:
		utils.FormatRelations(v)
	}
	return nil
}

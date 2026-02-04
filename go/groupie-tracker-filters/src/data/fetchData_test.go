package data

import (
	"encoding/json"
	"groupie-tracker/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// createMockArtist creates a test artist with predefined data
// This helper function ensures consistency across tests and reduces code duplication
func createMockArtist() models.Artist {
	return models.Artist{
		ID:           1,
		Image:        "http://example.com/image.jpg",
		Name:         "Test Artist",
		Members:      []string{"Member 1", "Member 2"},
		CreationDate: 2000,
		FirstAlbum:   "2001-01-01",
		Locations:    "http://localhost/api/locations",
		ConcertDates: "http://localhost/api/dates",
		Relations:    "http://localhost/api/relations",
	}
}

// TestFetchArtists tests the main FetchArtists function which retrieves artist data
// and their related information (locations, dates, relations)
func TestFetchArtists(t *testing.T) {
	// SETUP PHASE
	// ------------
	// Create a mock artist that will be returned by our fake API
	mockArtists := []models.Artist{createMockArtist()}

	// Create a test server that will respond with our mock artists data
	// This replaces the real API endpoint with a controlled test environment
	artistServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockArtists)
	}))
	// Ensure the test server is cleaned up after the test
	defer artistServer.Close()

	// Create a mock server for locations data
	// This simulates the API endpoint that provides venue locations
	locationsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(models.LocationsData{
			Locations: []string{"New York", "London"},
		})
	}))
	defer locationsServer.Close()

	// Create a mock server for concert dates
	// This simulates the API endpoint that provides performance dates
	datesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(models.DatesData{
			Dates: []string{"2024-01-01", "2024-02-01"},
		})
	}))
	defer datesServer.Close()

	// Create a mock server for relations data
	// This simulates the API endpoint that provides the relationship between locations and dates
	relationsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(models.RelationsData{
			DatesLocations: map[string][]string{
				"New York": {"2024-01-01"},
				"London":   {"2024-02-01"},
			},
		})
	}))
	defer relationsServer.Close()

	// CONFIGURATION PHASE
	// ------------------
	// Set the environment variable to point to our test server instead of the real API
	os.Setenv("myData", artistServer.URL)

	// Update the mock artist's URLs to point to our test servers
	mockArtists[0].Locations = locationsServer.URL
	mockArtists[0].ConcertDates = datesServer.URL
	mockArtists[0].Relations = relationsServer.URL

	// TEST EXECUTION PHASE
	// -------------------
	// Call the function we're testing
	artists, err := FetchArtists()

	// ASSERTION PHASE
	// --------------
	// Verify no errors occurred during the fetch
	if err != nil {
		t.Errorf("FetchArtists() returned unexpected error: %v", err)
	}

	// Verify we got the expected number of artists
	if len(artists) != 1 {
		t.Errorf("Expected 1 artist, got %d", len(artists))
	}

	// Verify the artist data was correctly parsed
	if artists[0].Name != "Test Artist" {
		t.Errorf("Expected artist name 'Test Artist', got '%s'", artists[0].Name)
	}
}

// TestFetchArtistsError tests error handling scenarios in the FetchArtists function
// This ensures the function handles API failures gracefully
func TestFetchArtistsError(t *testing.T) {
	// Test Case 1: Invalid URL
	// -----------------------
	// Set an invalid URL to simulate a network error
	os.Setenv("myData", "http://invalid-url")

	// Attempt to fetch artists with invalid URL
	_, err := FetchArtists()

	// Verify that an error was returned
	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}

	// Test Case 2: Invalid JSON Response
	// --------------------------------
	// Create a server that returns invalid JSON data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// Configure the test to use our server that returns invalid JSON
	os.Setenv("myData", server.URL)

	// Attempt to fetch artists with invalid JSON response
	_, err = FetchArtists()

	// Verify that an error was returned
	if err == nil {
		t.Error("Expected error with invalid JSON, got nil")
	}
}

// TestFetchRelatedDataAsync tests the concurrent fetching of related artist data
// This ensures that the async fetching of locations, dates, and relations works correctly
func TestFetchRelatedDataAsync(t *testing.T) {
	// SETUP PHASE
	// -----------
	artist := createMockArtist()

	// Create mock handlers for each type of data
	// Each handler returns test data for its respective endpoint
	handlers := map[string]http.HandlerFunc{
		"locations": func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(models.LocationsData{
				Locations: []string{"Test Location"},
			})
		},
		"dates": func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(models.DatesData{
				Dates: []string{"2024-01-01"},
			})
		},
		"relations": func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(models.RelationsData{
				DatesLocations: map[string][]string{
					"Test Location": {"2024-01-01"},
				},
			})
		},
	}

	// Create test servers for each endpoint
	for endpoint, handler := range handlers {
		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		// Update the artist's URLs to point to our test servers
		switch endpoint {
		case "locations":
			artist.Locations = server.URL
		case "dates":
			artist.ConcertDates = server.URL
		case "relations":
			artist.Relations = server.URL
		}
	}

	// TEST EXECUTION PHASE
	// -------------------
	// Fetch all related data concurrently
	err := fetchRelatedDataAsync(&artist)

	// ASSERTION PHASE
	// --------------
	// Verify no errors occurred during the fetch
	if err != nil {
		t.Errorf("fetchRelatedDataAsync() returned unexpected error: %v", err)
	}

	// Verify all data was populated correctly
	if len(artist.LocationsData.Locations) == 0 {
		t.Error("Locations data was not populated")
	}
	if len(artist.ConcertDatesData.Dates) == 0 {
		t.Error("Concert dates data was not populated")
	}
	if len(artist.RelationsData.DatesLocations) == 0 {
		t.Error("Relations data was not populated")
	}
}

// TestFetchAndUnmarshal tests the helper function that handles HTTP requests and JSON parsing
// This ensures the base functionality of making requests and parsing responses works correctly
func TestFetchAndUnmarshal(t *testing.T) {
	// SETUP & SUCCESSFUL CASE
	// ----------------------
	// Create test data that we expect to receive
	mockData := models.LocationsData{Locations: []string{"Test Location"}}

	// Create a test server that returns our mock data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	// TEST EXECUTION - Successful case
	var result models.LocationsData
	err := fetchAndUnmarshal(server.URL, &result)

	// Verify successful case
	if err != nil {
		t.Errorf("fetchAndUnmarshal() returned unexpected error: %v", err)
	}
	if len(result.Locations) != len(mockData.Locations) {
		t.Errorf("Expected %d locations, got %d", len(mockData.Locations), len(result.Locations))
	}

	// ERROR CASES
	// -----------
	// Define different error scenarios to test
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "Invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Return malformed JSON to test error handling
				w.Write([]byte("invalid json"))
			},
		},
		{
			name: "Server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Return a 500 error to test error handling
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
	}

	// Run each error test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server for this specific error case
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer server.Close()

			// Attempt to fetch and unmarshal, expecting an error
			var result models.LocationsData
			err := fetchAndUnmarshal(server.URL, &result)

			// Verify that an error was returned
			if err == nil {
				t.Errorf("Expected error for case %s, got nil", tt.name)
			}
		})
	}
}

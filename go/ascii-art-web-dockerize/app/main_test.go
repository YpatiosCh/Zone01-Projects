package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"ascii-art-web-stylize/tools"
)

// TestHomeHandler tests the home page handler for GET requests
func TestHomeHandler(t *testing.T) {
    req, err := http.NewRequest(http.MethodGet, "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(tools.HomeHandler)
    handler.ServeHTTP(rr, req)

    // Check that the status code is 200
    if status := rr.Code; status != http.StatusOK {
        t.Fatalf("Expected status code %d, got %d", http.StatusOK, status)
    }

    // Check for key elements in the HTML response
    body := rr.Body.String()

    // Verify the title tag
    if !strings.Contains(body, "<title>ASCII Art Generator</title>") {
        t.Errorf("Expected title not found in homeHandler response")
    }

    // Verify the form action for ASCII art generation
    if !strings.Contains(body, `form action="/ascii-art" method="post"`) {
        t.Errorf("Expected form action not found in homeHandler response")
    }

    // Verify the presence of a textarea for input text
    if !strings.Contains(body, `<textarea id="text" name="text" required>`) {
        t.Errorf("Expected textarea for text input not found in homeHandler response")
    }

    // Verify radio buttons for banner selection
    if !strings.Contains(body, `input type="radio" id="standard" name="banner" value="standard" required`) ||
       !strings.Contains(body, `input type="radio" id="shadow" name="banner" value="shadow"`) ||
       !strings.Contains(body, `input type="radio" id="thinkertoy" name="banner" value="thinkertoy"`) {
        t.Errorf("Expected banner radio buttons not found in homeHandler response")
    }

    // Check for the submit button
    if !strings.Contains(body, `<button type="submit">Generate</button>`) {
        t.Errorf("Expected submit button not found in homeHandler response")
    }
}


// TestAsciiArtHandler tests the ASCII art generation handler with valid input
func TestAsciiArtHandler(t *testing.T) {
	form := "text=Hello&banner=standard"
	req, err := http.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tools.AsciiArtHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("asciiArtHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "H") {
		t.Errorf("asciiArtHandler did not generate ASCII art correctly")
	}
}

// TestAsciiArtHandlerBadRequest tests ASCII art generation handler with missing form data
func TestAsciiArtHandlerBadRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/ascii-art", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tools.AsciiArtHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("asciiArtHandler returned wrong status code for bad request: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestApiAsciiArtHandler tests the API ASCII art handler with JSON input
func TestApiAsciiArtHandler(t *testing.T) {
	requestBody := `{"text": "Hello", "banner": "standard"}`
	req, err := http.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tools.ApiAsciiArtHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var response struct {
		AsciiArt []string `json:"ascii_art"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	// Log the actual ASCII art output for debugging
	for i, line := range response.AsciiArt {
		t.Logf("ASCII Art Line %d: %s", i, line)
	}

	// Modify expected ASCII art to match the exact output format
	expectedStart := " _    _" // adjust based on actual output
	if len(response.AsciiArt) == 0 || !strings.HasPrefix(response.AsciiArt[0], expectedStart) {
		t.Errorf("apiAsciiArtHandler did not generate ASCII art correctly, expected start with: %v", expectedStart)
	}
}

// TestPrintAsciiArt with additional checks on ASCII formatting
func TestPrintAsciiArt(t *testing.T) {
	bannerPath := "banners/standard.txt"
	expectedArtStart := " _    _" // Expected start of ASCII art for "H" in standard font

	result, err := tools.PrintAsciiArt("Hello", bannerPath)
	if err != nil {
		t.Fatalf("printAsciiArt returned error: %v", err)
	}

	// Log ASCII art for debugging
	for i, line := range result {
		t.Logf("Line %d: %s", i, line)
	}

	// Check that ASCII art output starts as expected
	if len(result) == 0 || !strings.HasPrefix(result[0], expectedArtStart) {
		t.Errorf("printAsciiArt did not generate ASCII art correctly, expected start with: %v", expectedArtStart)
	}
}

// TestGetAsciiArtForLetter with refined expected output check
func TestGetAsciiArtForLetter(t *testing.T) {
	file, err := os.Open("banners/standard.txt")
	if err != nil {
		t.Fatal("could not open banner file for testing:", err)
	}
	defer file.Close()

	letterArt, err := tools.GetAsciiArtForLetter(file, 'A')
	if err != nil {
		t.Fatalf("getAsciiArtForLetter returned error: %v", err)
	}

	// Expected ASCII pattern start for the letter 'A'
	expectedA := "    /\\    " // Adjust based on your ASCII font file content

	// Log each line of letter 'A' for debugging
	for i, line := range letterArt {
		t.Logf("Letter 'A' Line %d: %s", i, line)
	}

	// Check if the generated art for 'A' matches the expected start
	if len(letterArt) != 8 || !strings.HasPrefix(letterArt[1], expectedA) {
		t.Errorf("getAsciiArtForLetter did not generate correct ASCII art for 'A', expected line to start with: %v", expectedA)
	}
}

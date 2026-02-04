package handlers

import (
	"groupie-tracker/models"
	"html"
	"net/http"
	"strconv"
	"strings"
)

// SearchSuggestion defines the structure for autocomplete suggestions.
type SearchSuggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// SearchResult represents grouped search results by category.
type SearchResult struct {
	Artists       []ArtistMatch     `json:"artists"`
	Members       []MemberMatch     `json:"members"`
	Locations     []LocationMatch   `json:"locations"`
	CreationDates []CreationMatch   `json:"creationDates"`
	FirstAlbums   []FirstAlbumMatch `json:"firstAlbums"`
}

// ArtistMatch links an artist with a matching name.
type ArtistMatch struct {
	Artist models.Artist `json:"artist"`
	Name   string        `json:"name"`
}

// MemberMatch links an artist with a matching member.
type MemberMatch struct {
	Artist models.Artist `json:"artist"`
	Member string        `json:"member"`
}

// LocationMatch links an artist with a matching location.
type LocationMatch struct {
	Artist   models.Artist `json:"artist"`
	Location string        `json:"location"`
}

// CreationMatch links an artist with a matching creation date.
type CreationMatch struct {
	Artist       models.Artist `json:"artist"`
	CreationDate int           `json:"creationDate"`
}

// FirstAlbumMatch links an artist with a matching first album.
type FirstAlbumMatch struct {
	Artist     models.Artist `json:"artist"`
	FirstAlbum string        `json:"firstAlbum"`
}

// SearchHandler handles search queries and renders search results.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the search query from the URL parameters.
	rawQuery := r.URL.Query().Get("q")
	query := html.EscapeString(strings.TrimSpace(rawQuery))
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Perform a grouped search on the artist data.
	results := searchArtistsGrouped(Artists, query)

	// Prepare the data structure for rendering the search results.
	data := struct {
		Results SearchResult
		Query   string
	}{
		Results: results,
		Query:   query,
	}

	// Render the search results using the "search.html" template.
	err := Templates.ExecuteTemplate(w, "search.html", data)
	if err != nil {
		RenderErrorTemplate(w, http.StatusInternalServerError, "Error rendering template")
		return
	}
}

// searchArtistsGrouped performs a search on artist data and groups results by category.
func searchArtistsGrouped(artists []models.Artist, query string) SearchResult {
	query = strings.TrimSpace(strings.ToLower(query))
	results := SearchResult{}

	// Track artists that have already been matched by name.
	matchedArtists := make(map[string]bool) // Track by artist name instead of ID.

	// First pass - find direct artist name matches.
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), query) {
			results.Artists = append(results.Artists, ArtistMatch{
				Artist: artist,
				Name:   artist.Name,
			})
			matchedArtists[strings.ToLower(artist.Name)] = true
		}
	}

	// Second pass - find other matches.
	for _, artist := range artists {
		// Check members for matches.
		for _, member := range artist.Members {
			if strings.HasPrefix(strings.ToLower(member), query) {
				// Skip if this is a member match for their own artist entry.
				if matchedArtists[strings.ToLower(artist.Name)] &&
					strings.Contains(strings.ToLower(artist.Name), strings.ToLower(member)) {
					continue
				}

				results.Members = append(results.Members, MemberMatch{
					Artist: artist,
					Member: member,
				})
			}
		}

		// Check locations for matches.
		for _, location := range artist.LocationsData.Locations {
			if strings.Contains(strings.ToLower(location), query) {
				results.Locations = append(results.Locations, LocationMatch{
					Artist:   artist,
					Location: location,
				})
			}
		}

		// Check creation dates if the query is a 4-digit year.
		// Check creation dates for partial or exact matches
		if len(query) <= 4 { // Year format only
			creationDateStr := strconv.Itoa(artist.CreationDate)
			if strings.Contains(creationDateStr, query) {
				results.CreationDates = append(results.CreationDates, CreationMatch{
					Artist:       artist,
					CreationDate: artist.CreationDate,
				})

			}
		}

		// Check first album matches.
		if !strings.Contains(query, "-") && len(query) == 4 {
			// Skip year-only queries for albums.
			continue
		}

		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			results.FirstAlbums = append(results.FirstAlbums, FirstAlbumMatch{
				Artist:     artist,
				FirstAlbum: artist.FirstAlbum,
			})
		}
	}

	return results
}

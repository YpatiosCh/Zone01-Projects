package models

type Artist struct {
	ID               int           `json:"id"`
	Image            string        `json:"image"`
	Name             string        `json:"name"`
	Members          []string      `json:"members"`
	CreationDate     int           `json:"creationDate"`
	FirstAlbum       string        `json:"firstAlbum"`
	Locations        string        `json:"locations"`    // URL of locations
	ConcertDates     string        `json:"concertDates"` // URL of dates
	Relations        string        `json:"relations"`    // URL of relations
	LocationsData    LocationsData `json:"-"`            // results from URL of locations
	ConcertDatesData DatesData     `json:"-"`            // results from URL of dates
	RelationsData    RelationsData `json:"-"`            // results from URL of relation
}

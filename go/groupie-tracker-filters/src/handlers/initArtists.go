package handlers

import (
	"groupie-tracker/data"
	"groupie-tracker/models"
)

var Artists []models.Artist

func InitializeArtists() error {
	var err error
	Artists, err = data.FetchArtists()
	if err != nil {
		return err
	}
	return nil
}

package parser

import (
	"errors"
	"fmt"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/utils/validation"
	"net/http"
	"strings"
)

func ParseValuesToCreatePost(r *http.Request, data models.CreatePostData, categories []models.Category, config *config.AppConfig) (models.CreatePostData, int, error) {
	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDs := r.Form["categories"]

	// Save form data for re-displaying on error
	data.FormData["title"] = title
	data.FormData["content"] = content
	data.FormCategories = categoryIDs

	// Parse multipart form (for file uploads)
	if err := r.ParseMultipartForm(config.MaxImageSize + 10*1024*1024); err != nil {
		return data, http.StatusInternalServerError, errors.New("Error processing form. Please try again")
	}

	// Validate the post
	errorStruct := validation.ValidatePost(title, content, categoryIDs, categories)
	if errorStruct.Error() != nil {
		validationErr := strings.Join(errorStruct.Error(), " & ")
		return data, http.StatusBadRequest, errors.New(validationErr)
	}
	return data, http.StatusOK, nil
}

func ParseValuesToEdit(r *http.Request, data models.EditPostData, categories []models.Category, config *config.AppConfig) (bool, int, models.EditPostData, error) {
	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDs := r.Form["categories"]
	removeImage := r.FormValue("remove_image") == "on" // checkbox value
	// Save form data for re-displaying on error
	data.FormData["title"] = title
	data.FormData["content"] = content
	data.FormCategories = categoryIDs
	fmt.Println(data.FormCategories)

	// Parse multipart form (for file uploads)
	if err := r.ParseMultipartForm(config.MaxImageSize + 10*1024*1024); err != nil {
		return false, http.StatusInternalServerError, data, fmt.Errorf("Error processing form. Please try again")
	}

	// Validate the post
	errorStruct := validation.ValidatePost(title, content, categoryIDs, categories)
	if errorStruct.Error() != nil {
		validationErr := strings.Join(errorStruct.Error(), " & ")
		return false, http.StatusBadRequest, data, errors.New(validationErr)
	}
	return removeImage, http.StatusOK, data, nil
}

func ParseValuesToEditComment(r *http.Request, data models.EditCommentData) (string, error) {
	content := r.FormValue("content")

	errorStruct := validation.ValidateComment(content)
	if errorStruct.Error() != nil {
		validationErr := strings.Join(errorStruct.Error(), " & ")
		return "", errors.New(validationErr)
	}
	return content, nil
}

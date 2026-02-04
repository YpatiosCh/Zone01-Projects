package validation

import (
	"fmt"
	"forum/internal/models"
	"html"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

// PostValidationError represents a validation error for post creation
type ValidationError struct {
	ErrorSlice []string
}

// Error implements the error interface
func (e ValidationError) Error() []string {
	return e.ErrorSlice
}

// SanitizeInput cleans user input to prevent XSS attacks
func SanitizeInput(input string) string {
	// Escape HTML special characters
	sanitized := html.EscapeString(input)
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

// ValidatePost validates post data before saving to the database
func ValidatePost(title, content string, categoryIDs []string, existingCategories []models.Category) *ValidationError {
	// Check if title is empty
	var errorStruct ValidationError
	if title == "" {
		msg := "Title is required"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// Check if title is too long
	if len(title) > 100 {
		msg := "Title is too long (maximum 100 characters)"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// Check if content is empty
	if content == "" {
		msg := "Content is required"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// Check if content is too long
	if len(content) > 10000 {
		msg := "Content is too long (maximum 10000 characters)"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// Check if at least one category is selected
	if len(categoryIDs) == 0 {
		msg := "Please select at least one category"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// Check if more than 3 categories are selected
	if len(categoryIDs) > 3 {
		msg := "You can choose maximum 3 categories"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// All Categories map
	validCategorySet := make(map[string]struct{})
	for _, existingCategory := range existingCategories {
		validCategorySet[existingCategory.ID] = struct{}{}
	}
	// Check if Submited categories ar contained in map categories
	var validCategories []string
	for _, id := range categoryIDs {
		if _, exists := validCategorySet[id]; exists {
			validCategories = append(validCategories, id)
		}
	}
	// Check if all submitted categories were valid
	if len(validCategories) != len(categoryIDs) {
		msg := "Valid Category Found"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	// All validation passed
	return &errorStruct
}

// ValidateComment validates a comment's content
func ValidateComment(content string) *ValidationError {
	var errorStruct ValidationError
	// Validate content
	if len(content) == 0 {
		msg := "Comment is required"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	if len(content) < 3 {
		msg := "Comment must be at least 3 characters"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	if len(content) > 1000 {
		msg := "Comment must be less than 1000 characters"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	return &errorStruct
}

func ValidateRegistration(username, email, password string) *ValidationError {

	var errorStruct ValidationError

	if strings.TrimSpace(username) == "" {
		msg := "Username cannot be empty"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	if len(username) < 3 || len(username) > 20 {
		msg := "Username has to be bigger than 3 and less than 20 characters"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	reUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !reUsername.MatchString(username) {
		msg := "User name can have capital ,lower characters, number or '_'"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	reEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !reEmail.MatchString(email) {
		msg := "Incorrect mail address"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	if len(password) < 8 || len(password) > 20 {
		msg := "Password must have 8-20 characters"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	hasLowerLetter := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasCapitalLetter := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	for _, letter := range password {
		if letter < 33 || letter > 126 {
			msg := "Password must have only latin characters or numbers or symbols"
			errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
			break
		}
	}

	if !hasLowerLetter || !hasCapitalLetter || !hasNumber {
		msg := "Password must have at least 1 capital 1 lower letter and 1 number"
		errorStruct.ErrorSlice = append(errorStruct.ErrorSlice, msg)
	}

	return &errorStruct
}

// IsValidImageType checks if the content type is a valid image type
func IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	return validTypes[contentType]
}

func ValidateImage(imageHeader *multipart.FileHeader, imageFile multipart.File, maxImageSize int64) (bool, int, error) {
	var hasImage bool

	if imageHeader != nil {
		// Image was uploaded
		defer imageFile.Close()

		// Validate image size
		if imageHeader.Size > maxImageSize {
			return false, http.StatusBadRequest, fmt.Errorf("Image is too large. Maximum size is 20MB.")
		}

		// Validate image type
		contentType := imageHeader.Header.Get("Content-Type")
		if !IsValidImageType(contentType) {
			return false, http.StatusBadRequest, fmt.Errorf("Invalid image format. Supported formats: JPEG, PNG, GIF.")
		}

		hasImage = true
	}
	return hasImage, http.StatusOK, nil
}

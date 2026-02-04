package validation

import (
	"html"
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
func ValidatePost(title, content string, categoryIDs []string) *ValidationError {
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

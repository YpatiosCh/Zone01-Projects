package image

import (
	"fmt"
	uuid "forum/pkg/UUID"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveImageLocally(uploadDir string, imageHeader *multipart.FileHeader, imageFile multipart.File) (string, error) {

	// Ensure upload directory exists
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("Error creating upload directory. Please try again.")
	}

	// Generate unique filename
	ext := filepath.Ext(imageHeader.Filename)
	filename := uuid.GenerateUUID() + ext
	filePath := filepath.Join(uploadDir, filename)

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Error creating file: %v\n", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, imageFile)
	if err != nil {
		return "", fmt.Errorf("Error copying file: %v\n", err)
	}
	fmt.Println("saved image")
	return filePath, nil
}

func RemoveImageLocaly(uploadDir, image string) error {
	// Check if directory exists again
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		return fmt.Errorf("Directory still doesn't exist after creation: %s\n", uploadDir)
	}

	// Attempt to delete the image
	err := os.Remove(image)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	fmt.Println("removed image")
	return nil
}

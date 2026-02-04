package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	GitHubClientID     string `json:"github_client_id"`
	GitHubClientSecret string `json:"github_client_secret"`
	AppURL             string `json:"app_url"`
	MaxImageSize       int64  `json:"max_image_size"`
	Port               string `json:"port"`
	UploadDir          string `json:"upload_dir"`
}

func LoadConfig() *AppConfig {
	// Try to load from JSON file first
	if config := loadFromJSON("config.json"); config != nil {
		return config
	}

	// Fallback to default values
	return &AppConfig{
		GoogleClientID:     "test_google_client_id",
		GoogleClientSecret: "test_google_client_secret",
		GitHubClientID:     "test_github_client_id",
		GitHubClientSecret: "test_github_client_secret",
		AppURL:             "http://localhost:8080",
		MaxImageSize:       20 * 1024 * 1024,
		Port:               ":8080",
		UploadDir:          "./static/files/uploads",
	}
}

// loadFromJSON attempts to read the OAuth configuration from a JSON file.
// If the file cannot be read or decoded, it returns nil.
func loadFromJSON(filename string) *AppConfig {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var config AppConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil
	}

	return &config
}

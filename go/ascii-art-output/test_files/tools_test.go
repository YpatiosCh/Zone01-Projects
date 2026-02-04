package main

import (
	"os"
	"testing"
	"ascii-art-output/tools"
)
// TestToFileOrNotToFile tests the ToFileOrNotToFile function.
func TestToFileOrNotToFile(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Test cases
		{"--output=file.txt", true},   // Valid output flag with .txt
		{"--output=wrongfile", false}, // Invalid output flag (wrong filename)
		{"someText", false},           // Regular text (no output flag)
		{"-flag", false},              // Invalid flag starting with "-" but not fullfilling the rest of the flag requirements
		{"--output=", false},          // Invalid output flag with missing filename
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := tools.ToFileOrNotToFile(test.input)
			if result != test.expected {
				t.Errorf("ToFileOrNotToFile(%v) = %v; want %v", test.input, result, test.expected)
			}
		})
	}
}
// TestWriteToFile tests the WriteToFile function.
func TestWriteToFile(t *testing.T) {
	// Test data
	linesToPrint := []string{"Hello", "World", "!"}
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Ensure the file is removed after the test
	defer os.Remove(tmpFile.Name())
	// Close the file immediately as we just need the name
	tmpFile.Close()
	// Call the WriteToFile function
	tools.WriteToFile(linesToPrint, tmpFile.Name())
	// Read the file content back to verify
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	// Expected output
	expectedOutput := "Hello\nWorld\n!\n"
	// Assert the file content is as expected
	if string(content) != expectedOutput {
		t.Errorf("WriteToFile() wrote %v; want %v", string(content), expectedOutput)
	}
}
// CreateOutput test will also test PrintLines from array and LocateLines function
// as they are combined together to create the CreateOutput function
func TestCreateOutput(t *testing.T) {
	// Ensure the standard.txt file exists for testing
	fontFilePath := "../fonts/standard.txt"
	if _, err := os.Stat(fontFilePath); os.IsNotExist(err) {
		t.Fatalf("Font file %s does not exist for testing", fontFilePath)
	}
	tests := []struct {
		name       string
		input      string
		font       string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:       "Valid input with standard font",
			input:      "HELLO",
			font:       fontFilePath,
			wantOutput: []string{" _    _   ______   _        _         ____   ", "| |  | | |  ____| | |      | |       / __ \\  ", "| |__| | | |__    | |      | |      | |  | | ", "|  __  | |  __|   | |      | |      | |  | | ", "| |  | | | |____  | |____  | |____  | |__| | ", "|_|  |_| |______| |______| |______|  \\____/  ", "                                             ", "                                             "}, // Adjust according to actual content of standard.txt
			wantErr:    false,
		},
		{
			name:       "Empty input with standard font",
			input:      "",
			font:       fontFilePath,
			wantOutput: []string{""}, // Expect no output for empty input
			wantErr:    false,
		},
		{
			name:       "Invalid input with standard font",
			input:      "INVALID",
			font:       fontFilePath,
			wantOutput: []string{" _____   _   _  __      __             _        _____   _____   ", "|_   _| | \\ | | \\ \\    / /     /\\     | |      |_   _| |  __ \\  ", "  | |   |  \\| |  \\ \\  / /     /  \\    | |        | |   | |  | | ", "  | |   | . ` |   \\ \\/ /     / /\\ \\   | |        | |   | |  | | ", " _| |_  | |\\  |    \\  /     / ____ \\  | |____   _| |_  | |__| | ", "|_____| |_| \\_|     \\/     /_/    \\_\\ |______| |_____| |_____/  ", "                                                                ", "                                                                "}, // Adjust according to actual content of standard.txt
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := tools.CreateOutput(tt.input, tt.font)
			// Check if the result matches the expected output
			if len(gotOutput) != len(tt.wantOutput) {
				t.Errorf("CreateOutput() = %v; want %v", gotOutput, tt.wantOutput)
				return
			}
			// If we expect a specific output, check the contents
			for i, gotLine := range gotOutput {
				if gotLine != tt.wantOutput[i] {
					t.Errorf("CreateOutput() line %d = %v; want %v", i, gotLine, tt.wantOutput[i])
				}
			}
		})
	}
}

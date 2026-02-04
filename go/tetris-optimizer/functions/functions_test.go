// functions_test.go
package functions

import (
	"os"
	"testing"
)

func TestValidateAndCreateTetromino(t *testing.T) {
	tests := []struct {
		name      string
		shape     []string
		letter    rune
		wantError bool
	}{
		{
			name: "Valid I tetromino",
			shape: []string{
				"#...",
				"#...",
				"#...",
				"#...",
			},
			letter:    'A',
			wantError: false,
		},
		{
			name: "Valid square tetromino",
			shape: []string{
				"##..",
				"##..",
				"....",
				"....",
			},
			letter:    'B',
			wantError: false,
		},
		{
			name: "Invalid - not connected",
			shape: []string{
				"#...",
				"#...",
				"....",
				"##..",
			},
			letter:    'C',
			wantError: true,
		},
		{
			name: "Invalid - wrong number of blocks",
			shape: []string{
				"##..",
				"#...",
				"#...",
				"#...",
			},
			letter:    'D',
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateAndCreateTetromino(tt.shape, tt.letter)
			if (err != nil) != tt.wantError {
				t.Errorf("validateAndCreateTetromino() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestParseInput(t *testing.T) {
	// Create a temporary test file
	content := `#...
#...
#...
#...

....
....
..##
..##
`
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test parsing
	tetrominoes, err := ParseInput(tmpfile.Name())
	if err != nil {
		t.Fatalf("parseInput() error = %v", err)
	}

	if len(tetrominoes) != 2 {
		t.Errorf("Expected 2 tetrominoes, got %d", len(tetrominoes))
	}
}

func TestInitialBoardSize(t *testing.T) {
	tests := []struct {
		numTetrominoes int
		expectedSize   int
	}{
		{1, 2}, // 4 blocks = 2x2 minimum
		{2, 3}, // 8 blocks = 3x3 minimum
		{4, 4}, // 16 blocks = 4x4 minimum
		{8, 6}, // 32 blocks = 6x6 minimum
	}

	for _, tt := range tests {
		size := initialBoardSize(tt.numTetrominoes)
		if size != tt.expectedSize {
			t.Errorf("initialBoardSize(%d) = %d, want %d", tt.numTetrominoes, size, tt.expectedSize)
		}
	}
}

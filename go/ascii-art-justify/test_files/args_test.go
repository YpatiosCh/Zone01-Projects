package test_files

import (
	"ascii-art-justify/ascii"
	"flag"
	"os"
	"testing"
)

func TestHandleArgs(t *testing.T) {
	tests := []struct {
		args              []string
		expectedError     bool
		expectedSubstring string
		expectedFullInput string
		expectedColor     string
		expectedFont      string
		expectedOutput    string
		expectedAlign     string
	}{
		// Case 1: Valid input with substring only
		{
			args:              []string{"program", "Hello"},
			expectedError:     false,
			expectedSubstring: "Hello",
			expectedFullInput: "Hello",
			expectedColor:     "white",
			expectedFont:      "./fonts/standard.txt",
			expectedOutput:    "",
			expectedAlign:     "left",
		},
		// Case 2: Valid input with color flag and full input
		{
			args:              []string{"program", "--color=red", "Hi"},
			expectedError:     false,
			expectedSubstring: "Hi",
			expectedFullInput: "Hi",
			expectedColor:     "red",
			expectedFont:      "./fonts/standard.txt",
			expectedOutput:    "",
			expectedAlign:     "left",
		},
		// Case 3: Invalid alignment with output flag
		{
			args:          []string{"program", "--align=center", "--output=output.txt", "Hello"},
			expectedError: true,
		},
		// Case 4: Invalid color with output flag
		{
			args:          []string{"program", "--color=green", "--output=output.txt", "Test"},
			expectedError: true,
		},
		// Case 6: Valid input with substring and full input
		{
			args:              []string{"program", "sub", "full"},
			expectedError:     false,
			expectedSubstring: "sub",
			expectedFullInput: "full",
			expectedColor:     "white",
			expectedFont:      "./fonts/standard.txt",
			expectedOutput:    "",
			expectedAlign:     "left",
		},
		// Case 7: Missing arguments
		{
			args:          []string{"program"},
			expectedError: true,
		},
	}

	for _, test := range tests {
		// Reset and set up flags for each test case
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = test.args

		// Run the function and capture outputs
		substring, fullInput, color, font, output, align, err := ascii.HandleArgs()
		if (err != nil) != test.expectedError {
			t.Errorf("expected error: %v, got: %v, for args: %v", test.expectedError, err, test.args)
		}
		if err == nil {
			if substring != test.expectedSubstring {
				t.Errorf("expected substring: %v, got: %v, for args: %v", test.expectedSubstring, substring, test.args)
			}
			if fullInput != test.expectedFullInput {
				t.Errorf("expected fullInput: %v, got: %v, for args: %v", test.expectedFullInput, fullInput, test.args)
			}
			if color != test.expectedColor {
				t.Errorf("expected color: %v, got: %v, for args: %v", test.expectedColor, color, test.args)
			}
			if font != test.expectedFont {
				t.Errorf("expected font: %v, got: %v, for args: %v", test.expectedFont, font, test.args)
			}
			if output != test.expectedOutput {
				t.Errorf("expected output: %v, got: %v, for args: %v", test.expectedOutput, output, test.args)
			}
			if align != test.expectedAlign {
				t.Errorf("expected align: %v, got: %v, for args: %v", test.expectedAlign, align, test.args)
			}
		}
	}
}

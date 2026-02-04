package main

import (
	"fmt"
	"os"
	"text-modification/modification"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	text := string(data)

	finalText := totalFormat(text)
	// Write output file
	err = os.WriteFile(outputFile, []byte(finalText), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

// totalFormat formats the text using all the functions needed in order to get the final result of the project
func totalFormat(text string) string {
	createDecimals := modification.FormatDecimals(text)
	applyMethods := modification.ApplySingleMethods(createDecimals)
	applyMultMethods := modification.ApplyMultipleMethods(applyMethods)
	handlePunctuations := modification.HandlePunc(applyMultMethods)
	handleQuotes := modification.HandleQuotes(handlePunctuations)
	handleArticle := modification.HandleArticles(handleQuotes)
	handleArticle = modification.TrimTrailingSpaces(handleArticle)

	return handleArticle
}

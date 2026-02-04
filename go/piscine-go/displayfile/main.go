package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args[1:] // Retrieve command-line arguments excluding the program name

	length := len(arguments) // Determine the length of the arguments slice

	if length > 1 { // Check if the number of arguments is greater than 1
		fmt.Println("Too many arguments") // Print an error message for too many arguments
	} else if length == 0 { // Check if there are no arguments
		fmt.Println("File name missing") // Print an error message for missing file name
	} else if arguments[0] == "quest8.txt" { // Check if the argument is "quest8.txt"

		content, err := os.ReadFile(arguments[0]) // Read the content of the file named "quest8.txt"
		if err != nil {                           // Check if there was an error reading the file
			fmt.Println(err.Error()) // Print the error message
			return                   // Exit the program
		}
		fmt.Print(string(content)) // Print the content of the file as a string

	}
}

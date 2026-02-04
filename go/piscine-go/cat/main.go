package main

import (
	"io/ioutil"
	"os"

	"github.com/01-edu/z01"
)

func main() {
	if len(os.Args) == 1 {
		// If no arguments provided, read from standard input and print to standard output
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			printError("Error reading standard input: ", err)
			os.Exit(1)
		}
		printString(string(data))
	} else {
		// Iterate over each provided file and print its content
		for _, filename := range os.Args[1:] {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				printError("ERROR: "+err.Error(), nil)
				os.Exit(1)
			}
			printString(string(data))
		}
	}
}

func printError(msg string, err error) {
	for _, r := range msg {
		z01.PrintRune(r)
	}
	z01.PrintRune('\n')
	os.Exit(1)
}

func printString(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
}

package main

import (
	"ascii-art-fs/asciiart"
	"fmt"
	"os"
)

func main() {
	input := os.Args[1]

	font := "./fonts/standard.txt"

	if len(os.Args) > 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]\n\nEX: go run . something standard")
		return
	} else {
		text := os.Args[2]
		font = "./fonts/" + text + ".txt"
	}

	lines := asciiart.LocateLines(input)

	outputLine, err := asciiart.PrintLinesFromArray(font, lines)
	if err != nil {
		fmt.Print(err)
	}

	for _, line := range outputLine {
		fmt.Println(line)
	}
}

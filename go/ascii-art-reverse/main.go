package main

import (
	"reverse/args"
	"reverse/file"
)

func main() {
	filename, err := args.HandleArgs()
	if err != nil {
		args.PrintReverseUsage()
		// fmt.Println(err)
		return
	}

	asciiArt, err := file.ScanFile(filename)
	if err != nil {
		args.PrintReverseUsage()
		// fmt.Println(err)
		return
	}

	indexes := file.IndexOfEachLetter(asciiArt)
	file.PrintArtToAscii(indexes, asciiArt)
}

package main

import (
	"ascii-art-justify/ascii"
)

func main() {
	substring, fullInput, colorFlag, font, outputFlag, alignFlag, err := ascii.HandleArgs()
	if err != nil {
		ascii.PrintUsage()
		//fmt.Println(err)
		return
	}
	ascii.PrintFullAscii(substring, fullInput, colorFlag, outputFlag, font, alignFlag)
}

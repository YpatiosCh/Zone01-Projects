package main

import (
	"fmt"
	"math-skills/math"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file to read the data from")
	} else if len(os.Args) > 2 {
		fmt.Println("Please provide only the text file name to read the data")
	}

	fileName := os.Args[1]
	nums, err := math.ReadData(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	math.GetCalculations(nums)
}

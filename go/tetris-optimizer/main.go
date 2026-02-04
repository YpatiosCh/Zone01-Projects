package main

import (
	"fmt"
	"os"
	"tetris/functions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}

	tetrominoes, err := functions.ParseInput(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	solution := functions.SolvePuzzle(tetrominoes)
	if solution == nil {
		fmt.Println("ERROR")
		return
	}
	functions.PrintBoard(solution)
}

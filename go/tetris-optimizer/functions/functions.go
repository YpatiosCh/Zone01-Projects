package functions

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Tetromino struct {
	blocks []Point
	letter rune
}

type Board [][]rune

// ParseInput reads a file containing tetromino shapes, validates the format,
// and converts each valid shape into a Tetromino object. Returns a slice of
// Tetromino objects and an error if any shape is invalid or if there are no tetrominoes.
func ParseInput(filename string) ([]Tetromino, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tetrominoes []Tetromino
	scanner := bufio.NewScanner(file)
	currentShape := make([]string, 0, 4)
	letter := 'A'

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(currentShape) > 0 {
				tetromino, err := validateAndCreateTetromino(currentShape, letter)
				if err != nil {
					return nil, err
				}
				tetrominoes = append(tetrominoes, tetromino)
				letter++
				currentShape = make([]string, 0, 4)
			}
			continue
		}

		if len(line) != 4 || !isValidLine(line) {
			return nil, fmt.Errorf("invalid line format")
		}
		currentShape = append(currentShape, line)
	}

	// Handle the last tetromino
	if len(currentShape) > 0 {
		tetromino, err := validateAndCreateTetromino(currentShape, letter)
		if err != nil {
			return nil, err
		}
		tetrominoes = append(tetrominoes, tetromino)
	}

	if len(tetrominoes) == 0 {
		return nil, fmt.Errorf("no tetrominoes found")
	}

	return tetrominoes, nil
}

// isValidLine checks if a line in the input file consists only of '.' and '#' characters.
// Returns true if valid; false otherwise.
func isValidLine(line string) bool {
	for _, ch := range line {
		if ch != '.' && ch != '#' {
			return false
		}
	}
	return true
}

// validateAndCreateTetromino checks if the shape is a valid tetromino, normalizes it to a top-left corner,
// and returns it as a Tetromino object. Returns an error if the shape has an invalid height, block count,
// or disconnected blocks.
func validateAndCreateTetromino(shape []string, letter rune) (Tetromino, error) {
	if len(shape) != 4 {
		return Tetromino{}, fmt.Errorf("invalid tetromino height")
	}

	var blocks []Point
	blockCount := 0

	for y, line := range shape {
		for x, ch := range line {
			if ch == '#' {
				blocks = append(blocks, Point{x, y})
				blockCount++
			}
		}
	}

	if blockCount != 4 {
		return Tetromino{}, fmt.Errorf("invalid number of blocks")
	}

	if !isConnected(blocks) {
		return Tetromino{}, fmt.Errorf("blocks are not connected")
	}

	// Normalize the tetromino to top-left corner
	minX := blocks[0].x
	minY := blocks[0].y
	for _, p := range blocks {
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
	}

	for i := range blocks {
		blocks[i].x -= minX
		blocks[i].y -= minY
	}

	return Tetromino{blocks: blocks, letter: letter}, nil
}

// isConnected checks if all blocks in a tetromino shape are adjacent.
// Uses DFS to explore connectivity. Returns true if all blocks are connected; false otherwise.
func isConnected(blocks []Point) bool {
	if len(blocks) != 4 {
		return false
	}

	visited := make(map[Point]bool)
	visited[blocks[0]] = true
	stack := []Point{blocks[0]}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check all adjacent blocks
		neighbors := []Point{
			{current.x + 1, current.y},
			{current.x - 1, current.y},
			{current.x, current.y + 1},
			{current.x, current.y - 1},
		}

		for _, neighbor := range neighbors {
			for _, block := range blocks {
				if block.x == neighbor.x && block.y == neighbor.y && !visited[block] {
					visited[block] = true
					stack = append(stack, block)
				}
			}
		}
	}

	return len(visited) == 4
}

// SolvePuzzle attempts to place all tetrominoes on the smallest possible square board
// and increases board size until a solution is found. Returns the completed board.
func SolvePuzzle(tetrominoes []Tetromino) Board {
	size := initialBoardSize(len(tetrominoes))
	for {
		board := make(Board, size)
		for i := range board {
			board[i] = make([]rune, size)
			for j := range board[i] {
				board[i][j] = '.'
			}
		}

		if backtrack(board, tetrominoes, 0) {
			return board
		}
		size++
	}
}

// initialBoardSize calculates the smallest board dimension that can accommodate all tetrominoes.
// Returns the calculated size based on the number of tetrominoes.
func initialBoardSize(numTetrominoes int) int {
	// Each tetromino has 4 blocks, so we need at least sqrt(4 * n) size
	size := 2
	for size*size < numTetrominoes*4 {
		size++
	}
	return size
}

// backtrack attempts to place each tetromino on the board recursively.
// If all tetrominoes are placed successfully, returns true; otherwise, false.
func backtrack(board Board, tetrominoes []Tetromino, index int) bool {
	if index >= len(tetrominoes) {
		return true
	}

	size := len(board)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if canPlaceTetromino(board, tetrominoes[index], x, y) {
				placeTetromino(board, tetrominoes[index], x, y)
				if backtrack(board, tetrominoes, index+1) {
					return true
				}
				removeTetromino(board, tetrominoes[index], x, y)
			}
		}
	}
	return false
}

// canPlaceTetromino checks if a tetromino can be placed on the board at specified coordinates.
// Returns true if placement is possible; false otherwise.
func canPlaceTetromino(board Board, tetromino Tetromino, startX, startY int) bool {
	size := len(board)
	for _, block := range tetromino.blocks {
		x, y := startX+block.x, startY+block.y
		if x < 0 || x >= size || y < 0 || y >= size || board[y][x] != '.' {
			return false
		}
	}
	return true
}

// placeTetromino places a tetromino on the board at specified coordinates.
func placeTetromino(board Board, tetromino Tetromino, startX, startY int) {
	for _, block := range tetromino.blocks {
		board[startY+block.y][startX+block.x] = tetromino.letter
	}
}

// removeTetromino removes a tetromino from the board, clearing its blocks.
func removeTetromino(board Board, tetromino Tetromino, startX, startY int) {
	for _, block := range tetromino.blocks {
		board[startY+block.y][startX+block.x] = '.'
	}
}

// PrintBoard outputs the current state of the board to the console.
func PrintBoard(board Board) {
	for _, row := range board {
		fmt.Println(string(row))
	}
}

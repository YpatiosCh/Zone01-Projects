package piscine

import (
	"github.com/01-edu/z01"
)

const N = 8

var position = [N]int{}

func is_safe(queen_number, row_position int) bool {
	for i := 0; i < queen_number; i++ {
		other_row_pos := position[i]

		if other_row_pos == row_position || other_row_pos == row_position-(queen_number-i) || other_row_pos == row_position+(queen_number-i) {
			return false
		}
	}
	return true
}

func solve_puzzle(k int) {
	if k == N {
		for i := 0; i < N; i++ {
			z01.PrintRune(rune(position[i] + 49))
		}
		z01.PrintRune('\n')
	} else {
		for i := 0; i < N; i++ {
			if is_safe(k, i) {
				position[k] = i
				solve_puzzle(k + 1)
			}
		}
	}
}

func EightQueens() {
	solve_puzzle(0)
}

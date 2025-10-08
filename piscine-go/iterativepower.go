package piscine

func IterativePower(nb int, power int) int {
	if power < 0 {
		return 0
	}
	result := 1
	for ; power > 0; power-- {
		result *= nb
	}
	return result
}

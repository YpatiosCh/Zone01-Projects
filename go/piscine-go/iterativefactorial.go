package piscine

func IterativeFactorial(nb int) int {
	if nb < 0 {
		return 0
	}
	result := 1
	for ; nb > 1; nb-- {
		if result*nb < result {
			return 0
		}
		result *= nb
	}
	return result
}

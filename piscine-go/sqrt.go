package piscine

func Sqrt(nb int) int {
	if nb <= 0 {
		return 0
	}
	result := 1
	for result*result < nb && result*result >= result {
		result++
	}
	if result*result == nb {
		return result
	} else {
		return 0
	}
}

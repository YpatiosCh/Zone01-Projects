package piscine

func IsPrime(nb int) bool {
	if nb <= 1 || (nb > 2 && nb%2 == 0) || (nb > 3 && nb%3 == 0) {
		return false
	}
	for i := 2; i <= nb/2; i++ {
		if (nb/i)*i == nb {
			return false
		}
	}
	return true
}

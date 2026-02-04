package piscine

func LastRune(s string) rune {
	l := 'a'
	for _, letter := range s {
		l = letter
	}
	return l
}

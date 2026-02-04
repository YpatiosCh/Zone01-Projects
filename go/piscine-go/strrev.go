package piscine

func StrRev(s string) string {
	reverse := ""
	for _, a := range s {
		reverse = string(a) + reverse
	}
	return reverse
}

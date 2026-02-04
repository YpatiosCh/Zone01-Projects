package piscine

func JumpOver(str string) string {
	if len(str) < 3 {
		return "\n"
	}
	var result string
	for i := 0; i < len(str); i++ {
		if i%3 == 2 {
			result += string(str[i])
		}
	}
	return result + "\n"
}

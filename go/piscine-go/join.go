package piscine

func Join(strs []string, sep string) string {
	var result string
	for i, str := range strs {
		result += str
		if i < len(strs)-1 {
			result += sep
		}
	}
	return result
}

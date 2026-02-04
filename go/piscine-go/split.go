package piscine

func Split(s, sep string) []string {
	result := make([]string, 0)
	start := 0
	for i := range s {
		if len(s[i:]) >= len(sep) && s[i:i+len(sep)] == sep { // Check for full separator string
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + len(sep)
		}
	}
	// Add the last substring after the last separator (if it exists)
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

package piscine

func ShoppingSummaryCounter(str string) map[string]int {
	// Create a map to store the count of each item
	summary := make(map[string]int)
	// If the input string is empty, return an empty map
	if len(str) == 0 {
		summary[""] = 1
		return summary
	}
	// Initialize variables to keep track of item boundaries
	start := 0
	end := 0
	// Iterate through the input string
	for i := 0; i < len(str); i++ {
		// If the current character is a space or it's the end of the string
		if str[i] == ' ' || i == len(str)-1 {
			// Update the end index to the current position (or one past the end)
			end = i
			// Extract the item from the substring
			if i == len(str)-1 && str[i] != ' ' {
				end++
			}
			item := str[start:end]
			// Increment the count for the item in the summary map
			summary[item]++
			// Update the start index for the next item
			start = end + 1
		}
	}
	// If the last character is a space, account for the empty item at the end
	if str[len(str)-1] == ' ' {
		summary[""]++
	}
	return summary
}

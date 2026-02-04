package piscine

// Compact removes empty strings from a slice of strings and updates the slice accordingly.
// It returns the number of non-empty strings removed.
func Compact(ptr *[]string) int {
	add := 0         // Initialize a counter for non-empty strings
	var arr []string // Initialize an empty slice to store non-empty strings
	for _, v := range *ptr {
		// If the string is not empty, append it to the new slice and increment the counter
		if v != "" {
			arr = append(arr, v)
			add++
		}
	}
	*ptr = arr // Update the original slice with the non-empty strings
	return add // Return the number of non-empty strings removed
}

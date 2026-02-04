package piscine

func ShoppingListSort(slice []string) []string {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Compare the lengths of the strings
			if len(slice[j]) > len(slice[j+1]) {
				// Swap the strings if they are in the wrong order
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}

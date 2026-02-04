package piscine

func ConvertBase(s, t, p string) string {
	ln := 0                // Initialize variable 'ln' to store the length of string 's'.
	ln2 := 0               // Initialize variable 'ln2' to store the length of string 't'.
	ln3 := 0               // Initialize variable 'ln3' to store the length of string 'p'.
	ans := ""              // Initialize an empty string 'ans' to store the result.
	st_t := map[rune]int{} // Initialize a map 'st_t' to store the index of each character in string 't'.

	// Loop through each character in string 's' to determine its length.
	for c := range s {
		ln = c + 1
	}

	// Loop through each character in string 't'.
	for i, c := range t {
		// Store the index of each character in the 'st_t' map.
		st_t[c] = i
		ln2 = i + 1 // Update 'ln2' to store the length of string 't'.
	}

	// Loop through each character in string 'p' to determine its length.
	for c := range p {
		ln3 = c + 1
	}

	pw := 1  // Initialize variable 'pw' to store the power.
	cnt := 0 // Initialize variable 'cnt' to store the result of intermediate calculations.

	// Loop through the characters of string 's' in reverse order.
	for i := ln - 1; i >= 0; i-- {
		// Calculate the base conversion using the base 't' and 'p'.
		cnt = cnt + st_t[[]rune(s)[i]]*pw
		pw *= ln2
	}

	// Perform base conversion until 'cnt' becomes 0.
	for cnt != 0 {
		ans = string(p[cnt%ln3]) + ans // Append the converted character to the result string 'ans'.
		cnt /= ln3                     // Update 'cnt' for the next iteration.
	}

	// Return the final converted result.
	return ans
}

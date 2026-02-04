package piscine

func Enigma(a ***int, b *int, c *******int, d ****int) {
	// Temporarily store the value of a in a temporary variable
	tempA := ***a
	tempB := *b
	tempC := *******c
	tempD := ****d
	// Move the value of c into a
	*******c = tempA
	// Move the value of d into c
	****d = tempC
	// Move the value of b into d
	*b = tempD
	// Move the value of a into b
	***a = tempB
}

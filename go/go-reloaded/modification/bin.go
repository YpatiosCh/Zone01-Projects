package modification

import "strconv"

// BinToDecimal converts a binary string to a decimal integer
func BinToDecimal(binStr string) (int64, error) {
	decimal, err := strconv.ParseInt(binStr, 2, 64)
	if err != nil {
		return 0, err
	}
	return decimal, nil
}

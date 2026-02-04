package modification

import "strconv"

// HexToDecimal converts a hex string to a decimal integer
func HexToDecimal(hexStr string) (int64, error) {
	decimal, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return decimal, nil
}

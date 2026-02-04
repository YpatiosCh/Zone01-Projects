package modification

import (
	"strconv"
	"strings"

	"unicode"
)

// ApplySingleMethods applies the methods (up) (low) and (cap) to only the previous word.
func ApplySingleMethods(text string) string {
	result := strings.Split(text, " ")

	for i, word := range result {
		if word == "(up)" {
			result[i-1] = strings.ToUpper(result[i-1])
			result[i] = ""
		} else if word == "(low)" {
			result[i-1] = strings.ToLower(result[i-1])
			result[i] = ""
		} else if word == "(cap)" {
			result[i-1] = Capitalize(result[i-1])
			result[i] = ""
		}
	}
	modifiedResult := ConvertSliceToStr(result)
	finalResult := TrimTrailingSpaces(modifiedResult)
	return finalResult
}

// ApplyMultipleMethods applies the methods (up, <number>) (low, <number>) and (cap, <number>) to <number> previous words.
func ApplyMultipleMethods(text string) string {
	result := strings.Split(text, " ")

	for i, word := range result {
		if word == "(up," {
			count := ExtractCount(result[i+1])
			for j := 1; j <= count; j++ {
				result[i-j] = strings.ToUpper(result[i-j])
			}
			result[i] = ""
			result[i+1] = ""
		}
		if word == "(low," {
			count := ExtractCount(result[i+1])
			for j := 1; j <= count; j++ {
				result[i-j] = strings.ToLower(result[i-j])
			}
			result[i] = ""
			result[i+1] = ""
		}
		if word == "(cap," {
			count := ExtractCount(result[i+1])
			for j := 1; j <= count; j++ {
				result[i-j] = Capitalize(result[i-j])
			}
			result[i] = ""
			result[i+1] = ""
		}
	}
	modifiedResult := ConvertSliceToStr(result)
	finalResult := TrimTrailingSpaces(modifiedResult)
	return finalResult
}

// ExtractCount extracts the number from a string.
func ExtractCount(str string) int {
	digits := ""
	for _, char := range str {
		if unicode.IsDigit(char) {
			digits += string(char)
		}
	}
	if len(digits) > 0 {
		result, _ := strconv.Atoi(digits)
		return result
	}
	return 0
}

package modification

import (
	"strconv"
	"strings"
)

// FormatDecimals formats both (bin) and (hex) into decimals
func FormatDecimals(text string) string {
	result := strings.Split(text, " ")

	for i, word := range result {
		if word == "(hex)" {
			hexNum, _ := HexToDecimal(result[i-1])
			result[i-1] = strconv.Itoa(int(hexNum))
			result[i] = ""
		} else if word == "(bin)" {
			binNum, _ := BinToDecimal(result[i-1])
			result[i-1] = strconv.Itoa(int(binNum))
			result[i] = ""
		}
	}
	finalResult := ConvertSliceToStr(result)

	return finalResult
}

// ConvertSliceToStr converts slice into a string ignoring empty strings.
// Balances the whitespaces created after manipulating the slice.
func ConvertSliceToStr(slice []string) string {
	result := ""
	for i, word := range slice {
		if word == "" {
			continue
		} else if i != len(slice)-1 {
			result += word + " "
		} else {
			result += word
		}
	}
	return result
}

// TrimTrailingSpaces removes trailing whitespaces from the input string
func TrimTrailingSpaces(str string) string {
	return strings.TrimRight(str, " \t\n\r")
}

// Capitalize capitalizes the first letter of the word.
func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
}

// HandlePunc handles the punctuations in order to be in the correct place.
func HandlePunc(sentence string) string {
	var result string

	for i, r := range sentence {
		empty := ""
		space := " "
		if r == ' ' && !isPunc(rune(sentence[i+1])) {
			result += string(r)
			continue
		} else if r == ' ' && isPunc(rune(sentence[i+1])) {
			result += empty
		} else if i < len(sentence)-1 && isPunc(r) && !isPunc(rune(sentence[i+1])) && sentence[i+1] != ' ' {
			result += string(r)
			result += space
		} else {
			result += string(r)
		}
	}
	return result
}

// helper function to define whether a letter is a punction or not
func isPunc(r rune) bool {
	if r == '.' || r == '?' || r == ',' || r == '!' || r == ';' || r == ':' {
		return true
	}
	return false
}

// HandleQuotes handles the quotes in order to be in the correct place.
func HandleQuotes(s string) string {
	result := strings.Split(s, " ")
	seenQuote := false

	for i, word := range result {

		if strings.Contains(word, "'") && len(word) > 1 {
			seenQuote = true
		}

		if word == "'" && !seenQuote {
			result[i+1] = word + result[i+1]
			result[i] = ""
			seenQuote = true
		} else if word == "'" && seenQuote {
			result[i-1] = result[i-1] + word
			result[i] = ""
		} else {
			continue
		}
	}

	finalResult := ConvertSliceToStr(result)
	return finalResult
}

// HandleArticles handles the correction of the article 'a' -> 'an' and 'A' -> 'An' as needed.
func HandleArticles(text string) string {
	result := strings.Split(text, " ")

	for i, word := range result {
		if word == "a" && StartsWithVowelOrH(result[i+1]) {
			result[i] = "an"
		}
		if word == "A" && StartsWithVowelOrH(result[i+1]) {
			result[i] = "An"
		}
	}
	finalResult := ConvertSliceToStr(result)
	return finalResult
}

// StartsWithVowelOrH defines whether the first letter of a word starts with vowel or 'h'.
func StartsWithVowelOrH(s string) bool {
	return s[0] == 'a' || s[0] == 'e' || s[0] == '0' || s[0] == 'o' || s[0] == 'u' || s[0] == 'h'
}

package internals

import (
	"unicode"
)

type WordsResponse struct {
	Data []string `json:"data"`
}

func HasSpecialChar(input string) bool {
	for _, char := range input {
		if !unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

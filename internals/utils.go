package internals

import (
	"encoding/json"
	"io"
	"net/http"
	"unicode"
)

type WordsResponse struct {
	Data []string `json:"data"`
}

func GetWordsList() (*WordsResponse, error) {
	resp, err := http.Get("https://www.randomlists.com/data/words.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var words WordsResponse

	err = json.Unmarshal(body, &words)
	if err != nil {
		return nil, err
	}

	return &words, nil
}

func HasSpecialChar(input string) bool {
	for _, r := range input {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

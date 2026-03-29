package util

import (
	"strings"
)

func Tokenize(msg string) []string {
	msg = strings.ToLower(msg)
	var currToken strings.Builder
	var tokens []string
	for _, r := range msg {
		if IsSeparator(r) {
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
		} else {
			currToken.WriteRune(r)
		}
	}
	// add last token
	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}
	return tokens
}

func FilterStopWords(tokens []string) []string {
	var stopWords = map[string]struct{}{
		"a": {}, "an": {}, "the": {}, "and": {}, "or": {}, "in": {}, "on": {}, "at": {}, "is": {}, "to": {},
	}
	result := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if _, ok := stopWords[t]; !ok {
			result = append(result, t)
		}
	}
	return result
}

func IsSeparator(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', ',', ':', '=', '[', ']', '(', ')', '{', '}', '"', '\'', '/', '\\', '|', '-', '<', '>', ';', '.':
		return true
	}
	return false
}

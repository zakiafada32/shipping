// Package translation houses logic to take a given word and find it in a different language
package translation

import "strings"

type StaticService struct{}

func NewStaticService() *StaticService {
	return &StaticService{}
}

func (s *StaticService) Translate(word string, language string) string {
	word = sanitizeInput(word)         // <1>
	language = sanitizeInput(language) // <2>

	if word != "hello" {
		return ""
	}

	switch language {
	case "english":
		return "hello"
	case "finnish":
		return "hei"
	case "german":
		return "hallo"
	case "french":
		return "bonjour"
	default:
		return ""
	}
}

func sanitizeInput(w string) string { // <3>
	w = strings.ToLower(w)
	return strings.TrimSpace(w)
}

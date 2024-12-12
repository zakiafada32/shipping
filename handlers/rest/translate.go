// Package rest : Implement the HTTP handlers for the REST API.
package rest

import (
	"encoding/json"
	"net/http"
)

type Translator interface {
	Translate(word string, language string) string
}

type TranslateHandler struct {
	service Translator
}

func NewTranslateHandler(service Translator) *TranslateHandler {
	return &TranslateHandler{
		service: service,
	}
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

func (t *TranslateHandler) TranslateHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	language := r.URL.Query().Get("language")
	if language == "" {
		language = "english"
	}

	word := r.URL.Query().Get("word")

	translation := t.service.Translate(word, language)
	if translation == "" {
		w.WriteHeader(404)
		return
	}
	resp := Resp{
		Language:    language,
		Translation: translation,
	}
	if err := enc.Encode(resp); err != nil {
		panic("unable to encode response")
	}
}

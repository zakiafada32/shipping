package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zakiafada32/shipping-go/translation"
)

type Resp struct { // <2>
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	language := r.URL.Query().Get("language")
	if language == "" {
		language = "english"
	}

	word := r.URL.Query().Get("word")
	fmt.Printf("word: %s\n", word)
	translation := translation.Translate(word, language)
	if translation == "" {
		language = ""
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

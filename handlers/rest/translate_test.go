package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zakiafada32/shipping/handlers/rest"
)

type stubbedService struct{}

func (s *stubbedService) Translate(word string, language string) string {
	if word == "hello" {
		return "hallo"
	}
	return ""
}

func TestTranslateAPI(t *testing.T) {
	tt := []struct { // <1>
		Endpoint            string
		StatusCode          int
		ExpectedLanguage    string
		ExpectedTranslation string
	}{
		{
			Endpoint:            "?word=hello",
			StatusCode:          200,
			ExpectedLanguage:    "english",
			ExpectedTranslation: "hallo",
		},
		{
			Endpoint:            "?word=hello&language=german",
			StatusCode:          200,
			ExpectedLanguage:    "german",
			ExpectedTranslation: "hallo",
		},
		{
			Endpoint:   "?word=bye&language=japan",
			StatusCode: 404,
		},
	}

	h := rest.NewTranslateHandler(&stubbedService{})
	handler := http.HandlerFunc(h.TranslateHandler)

	for _, test := range tt {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.Endpoint, nil)

		handler.ServeHTTP(rr, req)

		if rr.Code != test.StatusCode {
			t.Errorf(`expected status %d but received %d`,
				test.StatusCode, rr.Code)
		}

		var resp rest.Resp
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)

		if resp.Language != test.ExpectedLanguage {
			t.Errorf(`expected language "%s" but received %s`,
				test.ExpectedLanguage, resp.Language)
		}

		if resp.Translation != test.ExpectedTranslation {
			t.Errorf(`expected Translation "%s" but received "%s"`,
				test.ExpectedTranslation, resp.Translation)
		}

	}
}

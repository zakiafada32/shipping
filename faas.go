// Package faas : Implement the HTTP handlers for the cloud function.
package faas

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/zakiafada32/shipping/handlers/rest"
	"github.com/zakiafada32/shipping/translation"
)

func init() {
	functions.HTTP("Translate", Translate)
}

func Translate(w http.ResponseWriter, r *http.Request) {
	translationService := translation.NewStaticService()
	translateHandler := rest.NewTranslateHandler(translationService)

	translateHandler.TranslateHandler(w, r)
}

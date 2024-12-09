package faas

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/zakiafada32/shipping-go/handlers/rest"
)

func init() {
	functions.HTTP("Translate", rest.TranslateHandler)
}

func Translate(w http.ResponseWriter, r *http.Request) {
	rest.TranslateHandler(w, r)
}

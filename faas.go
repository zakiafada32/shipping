package faas

import (
	"net/http"

	"github.com/zakiafada32/shipping-go/handlers/rest"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	rest.TranslateHandler(w, r)
}

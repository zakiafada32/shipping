package faas

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/zakiafada32/shipping-go/handlers/rest"
)

func init() {
	functions.HTTP("Translate", rest.TranslateHandler)
}

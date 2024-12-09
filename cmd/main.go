package main

import (
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/zakiafada32/shipping-go/handlers"
	"github.com/zakiafada32/shipping-go/handlers/rest"
)

func init() {
	functions.HTTP("translate", rest.TranslateHandler)
}

func main() {

	addr := ":8081"

	mux := http.NewServeMux()

	mux.HandleFunc("/translate", rest.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	log.Printf("listening on %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zakiafada32/shipping-go/handlers"
	"github.com/zakiafada32/shipping-go/handlers/rest"
	"github.com/zakiafada32/shipping-go/translation"
)

func main() {

	addr := ":8080"

	mux := http.NewServeMux()

	translationService := translation.NewStaticService()
	translateHandler := rest.NewTranslateHandler(translationService)
	mux.HandleFunc("/translate/hello", translateHandler.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck) // <1>

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s\n", addr)

	log.Fatal(server.ListenAndServe())
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zakiafada32/shipping-go/config"
	"github.com/zakiafada32/shipping-go/handlers"
	"github.com/zakiafada32/shipping-go/handlers/rest"
	"github.com/zakiafada32/shipping-go/translation"
)

func main() {

	cfg := config.LoadConfiguration()
	addr := cfg.Port

	mux := API(cfg)

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

func API(cfg config.Configuration) *http.ServeMux {

	mux := http.NewServeMux()

	var translationService rest.Translator
	translationService = translation.NewStaticService()
	if cfg.LegacyEndpoint != "" {
		log.Printf("creating external translation client: %s", cfg.LegacyEndpoint)
		client := translation.NewHelloClient(cfg.LegacyEndpoint)
		translationService = translation.NewRemoteService(client)
	}
	if cfg.DatabaseURL != "" {
		if cfg.DatabaseURL != "" {
			conn := translation.NewDatabaseService(cfg)
			conn.LoadData()
			translationService = conn
		}
	}
	translateHandler := rest.NewTranslateHandler(translationService)

	mux.HandleFunc("/translate", translateHandler.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	return mux
}

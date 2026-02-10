package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"price_checker/internal/features/price_tracker/repository"
	"price_checker/internal/features/price_tracker/service"
	"price_checker/internal/features/price_tracker/transport"
	"price_checker/internal/features/price_tracker/scraper"
)

func main() {

	repo := repository.NewStorage()

	htmlScraper := scraper.NewGoQueryScrapper()

	svc := service.NewPriceService(repo, htmlScraper)

	handler := transport.NewHandler(svc)

	// update prices of items
	svc.StartChecking(1 * time.Minute)

	// create router
	mux := http.NewServeMux()
	
	// set router
	mux.HandleFunc("POST /items", handler.AddItem)
	mux.HandleFunc("GET /items", handler.ListItems)

	log.Println("Starting server on :9091")

	if err := http.ListenAndServe(":9091", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed under request")
		} else {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}
}
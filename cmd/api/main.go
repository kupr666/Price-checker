package api

import (
	"log"
	"errors"
	"net/http"
	 
	"price_checker/internal/features/price_checker/repository"
	"price_checker/internal/features/price_checker/transport"
	"price_checker/internal/features/price_tracker/service"
)

func main() {

	repo := repository.NewStorage()

	svc := service.NewPriceService(repo)

	handler := transport.NewHandler(svc)

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
package main

import (
	"errors"
	// "log"
	"net/http"
	"time"

	"price_checker/internal/features/price_tracker/repository"
	"price_checker/internal/features/price_tracker/scraper"
	"price_checker/internal/features/price_tracker/service"
	"price_checker/internal/features/price_tracker/transport"
	"price_checker/internal/pkg/logger"
	"price_checker/internal/pkg/notifier"

	"go.uber.org/zap"
)

func main() {

	l, logFileClose, err := logger.NewLogger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	defer l.Sync()

	tgToken := ""
	tgChatID := ""


	repo := repository.NewStorage()
	tgNotifier := notifier.NewTelegramNotifier(tgToken, tgChatID)

	htmlScraper := scraper.NewGoQueryScrapper(l)

	svc := service.NewPriceService(repo, htmlScraper, l, tgNotifier)

	handler := transport.NewHandler(svc)

	// update prices of items
	svc.StartChecking(1 * time.Minute)

	// create router
	mux := http.NewServeMux()
	
	// set router
	mux.HandleFunc("POST /items", handler.AddItem)
	mux.HandleFunc("GET /items", handler.ListItems)

	l.Info("Starting server", zap.String("port", ":9091"))
	
	if err := http.ListenAndServe(":9091", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			l.Info("Server closed under request")
		} else {
			l.Error("Server forced to shutdown: %v", zap.Error(err))
		}
	}
}
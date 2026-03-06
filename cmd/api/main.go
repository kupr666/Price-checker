package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"price_checker/internal/features/price_tracker/repository"
	"price_checker/internal/features/price_tracker/scraper"
	"price_checker/internal/features/price_tracker/service"
	"price_checker/internal/features/price_tracker/transport"
	"price_checker/internal/pkg/db"
	"price_checker/internal/pkg/logger"
	"price_checker/internal/pkg/notifier"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	logLevel := os.Getenv("LOG_LEVEL")
	serverPort := os.Getenv("SERVER_PORT")


	l, logFileClose, err := logger.NewLogger(logLevel)
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	defer l.Sync()

	// global application context needed to graceful shotdown
	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	if err := db.RunMigrations(os.Getenv("CONNECTION_STRING")); err != nil {
		l.Fatal("failed to run migration", zap.Error(err))
	}

	pool, err := pgxpool.New(appCtx, os.Getenv("CONNECTION_STRING"))
	if err != nil {
		l.Fatal("failed to connect db", zap.Error(err))
	}
	defer pool.Close()

	repo := repository.NewPostgresStorage(pool)
	tgNotifier := notifier.NewTelegramNotifier(telegramToken, telegramChatID)
	htmlScraper := scraper.NewGoQueryScrapper(l)

	svc := service.NewPriceService(repo, htmlScraper, l, tgNotifier)
	handler := transport.NewHandler(svc, l)

	// update prices of items
	svc.StartChecking(appCtx, 1 * time.Minute)

	// create router
	mux := http.NewServeMux()
	// set up router
	mux.HandleFunc("POST /items", handler.AddItem)
	mux.HandleFunc("GET /items", handler.ListItems)
	mux.HandleFunc("DELETE /items/{id}", handler.DeleteItem)

	srv := &http.Server{
		Addr:	serverPort,
		Handler: mux,
	}

	// start the server in separte gorutine
	go func() {
		l.Info("Starting server", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	
	quit := make(chan os.Signal, 1)
	// when we press ctrl+c instead of killing the app instantly
	// this func sends message to quit channel
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// here we block and  wait this message
	<-quit
	l.Info("Shotdown signal received. Starting Graceful Shotdown")
	
	// cancel global background
	cancelApp()

	// give 10 seconds to close all jobs
	serverCtx, serverCancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer serverCancel()

	// finish app after 10 minutes
	if err := srv.Shutdown(serverCtx); err != nil {
		l.Error("Server forced to shutdown abnormally", zap.Error(err))
	} else {
		l.Info("Server exited gracefully")
	}
}
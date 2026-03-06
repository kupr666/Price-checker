package service

import (
	"fmt"
	// "log"
	"price_checker/internal/core/domains"
	"price_checker/internal/features/price_tracker/scraper"
	"sync"
	"time"
	"context"

	"go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
)

type Notifier interface {
	Notify(message string) error
}

type Repository interface {
	Add(context.Context, domains.Item) (domains.Item, error)
	Delete(context.Context, int64) error
	GetAll(context.Context) ([]domains.Item, error)
	UpdatePrice(context.Context, int64, float64) error
}

type PriceService struct {
	repo Repository
	scraper scraper.Scraper
	logger *zap.Logger
	notifier Notifier
}

func NewPriceService(repo Repository, sc scraper.Scraper, logger *zap.Logger, n Notifier) *PriceService {
	return &PriceService{
		repo: repo,
		scraper: sc,
		logger: logger,
		notifier: n,
	}
}

// start endless loop which checks price according to the interval
func (s *PriceService) StartChecking(ctx context.Context, interval time.Duration) {

	ticker := time.NewTicker(interval)

	go func () {
		defer ticker.Stop()
		for {
			select {
			case <- ticker.C:
				s.logger.Info("--- Starting background price check")
				s.CheckAllPrices(ctx)	
			case <- ctx.Done():
				s.logger.Info("Background worker received stop signal")
				return
			}	
		}
	}()

}

func (s *PriceService) CheckAllPrices(ctx context.Context) {

	// s.logger.Info("Starting background price check cycle")

	items, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve items from repository", zap.Error(err))
		return
	}

	// count of gorutines which work at the same time
	workerCount := 5

	// craete buffered channel with capacity = len(items) in order msg (jobs <- item) doesn't block
	jobs := make (chan domains.Item, len(items))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// worker hear channel till channel isn't closed
			// 5 gorutines immediately creates. The channel is empty -> workers block and wait items
			for item := range jobs {

				// s.logger.Debug("worker processing item",
				// 	zap.Int64("item_id", item.ID),
				// 	zap.String("item_url", item.URL),
				// )

				newPrice, err := s.scraper.FetchCurrentPrice(ctx, item.URL)
				if err != nil {
					s.logger.Error("scraping failed", zap.Error(err), zap.String("url_item", item.URL))
					continue
				}

				if err := s.repo.UpdatePrice(ctx, item.ID, newPrice); err != nil {
					s.logger.Error("failed to update price", zap.Error(err), zap.Int64("item_id", item.ID))
				}

				if newPrice <= item.TargetPrice {
					s.logger.Info("Target reached",
						zap.String("url", item.URL),
						zap.Float64("price", newPrice),
						zap.Float64("target", item.TargetPrice),
					)

					msg := fmt.Sprintf("Target reached!\nItem: %s\nNew Price: %.2f\nTarget: %.2f",
					item.URL, newPrice, item.TargetPrice)

					if err := s.notifier.Notify(msg); err != nil {
						// s.logger.Error("failed to send notification", zap.Error(err))
					}
				}
			}
		}(i)
	}

	for _, item := range items {
		// send item to a channel
		jobs <- item
	}

	// closing a channel sends a signal to all listeners (range jobs) that no more data will be sent
	// without this string workers would wait forever 
	close(jobs)

	// main gorutine musn't finish untill workers are done
	wg.Wait()

	s.logger.Info("Price check cycle completed", zap.Int("items_processed", len(items)))
}


func (s *PriceService) CreateItem(ctx context.Context, item domains.Item) (domains.Item, error) {

	if item.URL == "" || item.TargetPrice <= 0 {
		return domains.Item{}, fmt.Errorf("Invalid input: url and target price are required")
	}

	s.logger.Info("Creating item for url", zap.String("url", item.URL))

	currentPrice, err := s.scraper.FetchCurrentPrice(ctx, item.URL)
	if err != nil {
		return domains.Item{}, fmt.Errorf("couldn't fetch price: %w", err)
	}

		
	item.CurrentPrice = currentPrice	
	item.LastChecked = time.Now()	

	return s.repo.Add(ctx, item)
}

func (s *PriceService) ListItems(ctx context.Context) ([]domains.Item, error) {
	return s.repo.GetAll(ctx)
}

func (s *PriceService) DeleteItem(ctx context.Context, id int64) (error) {
	return s.repo.Delete(ctx, id)
}

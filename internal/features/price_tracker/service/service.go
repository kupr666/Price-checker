package service

import (
	"fmt"
	"log"
	"price_checker/internal/core/domains"
	"sync"
	"time"
	"price_checker/internal/features/price_tracker/scraper"

)

type Repository interface {
	Add(domains.Item) (domains.Item, error)
	Delete(int64) error
	GetAll() ([]domains.Item, error)
	UpdatePrice(int64, float64) error
}

type PriceService struct {
	repo Repository
	scraper scraper.Scraper
}

func NewPriceService(repo Repository, sc scraper.Scraper) *PriceService {
	return &PriceService{
		repo: repo,
		scraper: sc,
	}
}

// start endless loop which checks price according to the interval
func (s *PriceService) StartChecking(interval time.Duration) {

	ticker := time.NewTicker(interval)

	go func () {
		for range ticker.C {
			log.Println("--- Starting background price check")
			s.CheckAllPrices()	
		}
	}()

}

func (s *PriceService) CheckAllPrices() {
	items, err := s.repo.GetAll()
	if err != nil {
		log.Printf("Worker error: failed to get items: %v", err)
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
				newPrice, err := s.scraper.FetchCurrentPrice(item.URL)
				if err != nil {
					log.Printf("Error: %v", err)
					continue
				}
				s.repo.UpdatePrice(item.ID, newPrice)
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
}



func (s *PriceService) CreateItem(item domains.Item) (domains.Item, error) {

	if item.URL == "" || item.TargetPrice <= 0 {
		return domains.Item{}, fmt.Errorf("Invalid input: url and target price are required")
	}

	log.Printf("Creating item for url: %s", item.URL)


	currentPrice, err := s.scraper.FetchCurrentPrice(item.URL)
	if err != nil {
		return domains.Item{}, fmt.Errorf("couldn't fetch price: %w", err)
	}

		
	item.CurrentPrice = currentPrice	
	item.LastChecked = time.Now()	

	return s.repo.Add(item)
}

func (s *PriceService) ListItems() ([]domains.Item, error) {
	return s.repo.GetAll()
}


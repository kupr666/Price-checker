package repository

import (
	"price_checker/internal/core/domains"	
	"sync"
	"time"
	"context"
)

type Storage struct {
	mu 		sync.RWMutex
	items 	map[int64]domains.Item
	nextID 	int64
}

func NewStorage() *Storage {
	return &Storage{
		items: make(map[int64]domains.Item),
		nextID: 1,
	}
}

func (s *Storage) Add(ctx context.Context, item domains.Item) (domains.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if the context is already cancelled
	if err := ctx.Err(); err != nil {
		return domains.Item{}, err
	}

	for _, exists := range s.items {
		if exists.URL == item.URL {
			return domains.Item{}, domains.ErrUrlExists
		}
	}
	item.ID = s.nextID
	item.LastChecked = time.Now()
	s.items[item.ID] = item
	s.nextID++

	return item, nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return domains.ErrItemNotFound
	}

	delete(s.items, id)
	return nil
}

func (s *Storage) GetAll(ctx context.Context) ([]domains.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domains.Item, 0, len(s.items))

	for _, item := range s.items {
		items = append(items, item)
	}

	return items, nil
}

func (s *Storage) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return domains.ErrItemNotFound
	}

	item.CurrentPrice = newPrice
	item.LastChecked = time.Now()

	s.items[id] = item
	return nil
}
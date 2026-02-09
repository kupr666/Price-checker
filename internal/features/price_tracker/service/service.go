package service

import (
	"price_checker/internal/core/domains"
)

type Repository interface {
	Add(domains.Item) (domains.Item, error)
	Delete(int64) error
	GetAll() ([]domains.Item, error)
	UpdatePrice(int64, float64) error
}

type PriceService struct {
	repo Repository
}

func NewPriceService(repo Repository) *PriceService {
	return &PriceService{repo: repo}
}

func (s *PriceService) CreateItem(item domains.Item) (domains.Item, error) {
	return s.repo.Add(item)
}

func (s *PriceService) ListItems() ([]domains.Item, error) {
	return s.repo.GetAll()
}
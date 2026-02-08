package repository

import "price_checker/internal/core/domains"

type ItemStorage interface {
	Add(item domains.Item) (domains.Item, error)
	Delete(id int64) error
	GetAll() ([]domains.Item, error)
	UpdatePrice(id int64, newPrice float64) error
}
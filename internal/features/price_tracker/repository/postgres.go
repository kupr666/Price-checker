package repository

import (
	"time"
	"context"

	"price_checker/internal/core/domains"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage{
	return &PostgresStorage{
		pool: pool,
	}
}

func (s *PostgresStorage) Add(ctx context.Context, item domains.Item) (domains.Item, error) {

	item.LastChecked = time.Now()
	query := `
			INSERT INTO items (url, current_price, target_price, last_checked)
			VALUES ($1, $2, $3, $4)
			
	
	`

}

func (s *PostgresStorage) Delete(ctx context.Context, id int64) error {
	
}

func (s *PostgresStorage) GetAll (ctx context.Context) ([]domains.Item, error) {

}

func (s *PostgresStorage) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {

}
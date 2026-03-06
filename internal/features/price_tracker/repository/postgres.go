package repository

import (
	"fmt"
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
			RETURNING id
	`
	
	err := s.pool.QueryRow(ctx, query, item.URL, item.CurrentPrice, item.TargetPrice, item.LastChecked).Scan(&item.ID)
	if err != nil {
		return domains.Item{}, fmt.Errorf("add item: %w", err)
	}

	return item, nil
}

func (s *PostgresStorage) Delete(ctx context.Context, id int64) error {

	query := `DELETE FROM items WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delelte item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domains.ErrItemNotFound
	}

	return nil
}

func (s *PostgresStorage) GetAll(ctx context.Context) ([]domains.Item, error) {

	query := `SELECT id, url, current_price, target_price, last_checked FROM items`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all items: %w", err)
	}
	defer rows.Close()

	var items []domains.Item
	for rows.Next() {
		var item domains.Item
		if err := rows.Scan(&item.ID, &item.URL, &item.CurrentPrice,
			&item.TargetPrice, &item.LastChecked); err != nil {
				return nil, fmt.Errorf("scan item: %w", err)
			}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *PostgresStorage) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {

	query := `UPDATE items SET current_price = $1, last_checked = $2 WHERE id = $3`

	result, err := s.pool.Exec(ctx, query, newPrice, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update price: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domains.ErrItemNotFound
	}

	return nil
}
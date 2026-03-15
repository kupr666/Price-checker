package cache

import (
	"context"
	"fmt"
	// "os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type PriceCache interface {
	Get(ctx context.Context, url string) (float64, bool, error)
	Set(ctx context.Context, url string, price float64) error
}

type RedisCache struct {
	client *redis.Client
	ttl time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db, 				  // hardcoded (zero by default)
	})

	return &RedisCache{
		client: client, 
		ttl: ttl, 
	}
}

// find a cached price by URL
// Returns price, true, nil -> cache found
// Returns 0, false, nil -> key expired or never set
// Returns 0, false, error -> redis err
func (c *RedisCache) Get(ctx context.Context, url string) (float64, bool, error) {
	key := buildKey(url)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("redis get: %w", err)
	}

	price, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, false, fmt.Errorf("redis parse cached value %q: %w", val, err)
	}

	return price, true, nil
}

func (c *RedisCache) Set(ctx context.Context, url string, price float64) error {
	key := buildKey(url)

	// convert float into a string for redis
	val := strconv.FormatFloat(price, 'f', 2, 64)

	if err := c.client.Set(ctx, key, val, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func buildKey(url string) string {
	return "price:" + url 
}
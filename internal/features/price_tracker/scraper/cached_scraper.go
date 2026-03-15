package scraper

import (
	"context"

	"price_checker/internal/cache"

	"go.uber.org/zap"
)

type CachedSraper struct {
	scraper Scraper
	cache cache.PriceCache
	logger *zap.Logger
}

func NewCacheScraper(scraper Scraper, cache cache.PriceCache, logger *zap.Logger) *CachedSraper{
	return&CachedSraper{
		scraper: scraper,
		cache: cache,
		logger: logger,
	}
}

func (c *CachedSraper) FetchCurrentPrice(ctx context.Context, itemURL string) (float64, error) {

	found с.cache
}